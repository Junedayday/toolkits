package mysql

import (
	"context"
	"fmt"
	"net"

	vtmysql "vitess.io/vitess/go/mysql"
)

// BinLogReader implements reading bin log
type BinLogReader interface {
	ReadBinLog(ctx context.Context, serverID, offset uint32, filename string, flags uint16) (<-chan vtmysql.BinlogEvent, error)
}

// NewBinLogReader creates a bin log reader
// move from drvier.go : func (d MySQLDriver) Open(dsn string) (driver.Conn, error)
func NewBinLogReader(dsn string) (logger BinLogReader, err error) {
	// New mysqlConn
	mc := &mysqlConn{
		maxAllowedPacket: maxPacketSize,
		maxWriteSize:     maxPacketSize - 1,
		closech:          make(chan struct{}),
	}
	mc.cfg, err = ParseDSN(dsn)
	if err != nil {
		return nil, err
	}
	mc.parseTime = mc.cfg.ParseTime

	// Connect to Server
	dialsLock.RLock()
	dial, ok := dials[mc.cfg.Net]
	dialsLock.RUnlock()
	if ok {
		mc.netConn, err = dial(mc.cfg.Addr)
	} else {
		nd := net.Dialer{Timeout: mc.cfg.Timeout}
		mc.netConn, err = nd.Dial(mc.cfg.Net, mc.cfg.Addr)
	}
	if err != nil {
		return nil, err
	}

	// Enable TCP Keepalives on TCP connections
	if tc, ok := mc.netConn.(*net.TCPConn); ok {
		if err := tc.SetKeepAlive(true); err != nil {
			// Don't send COM_QUIT before handshake.
			mc.netConn.Close()
			mc.netConn = nil
			return nil, err
		}
	}

	// Call startWatcher for context support (From Go 1.8)
	mc.startWatcher()

	mc.buf = newBuffer(mc.netConn)

	// Set I/O timeouts
	mc.buf.timeout = mc.cfg.ReadTimeout
	mc.writeTimeout = mc.cfg.WriteTimeout

	// Reading Handshake Initialization Packet
	authData, plugin, err := mc.readHandshakePacket()
	if err != nil {
		mc.cleanup()
		return nil, err
	}
	if plugin == "" {
		plugin = defaultAuthPlugin
	}

	// Send Client Authentication Packet
	authResp, addNUL, err := mc.auth(authData, plugin)
	if err != nil {
		// try the default auth plugin, if using the requested plugin failed
		errLog.Print("could not use requested auth plugin '"+plugin+"': ", err.Error())
		plugin = defaultAuthPlugin
		authResp, addNUL, err = mc.auth(authData, plugin)
		if err != nil {
			mc.cleanup()
			return nil, err
		}
	}
	if err = mc.writeHandshakeResponsePacket(authResp, addNUL, plugin); err != nil {
		mc.cleanup()
		return nil, err
	}

	// Handle response to auth packet, switch methods if possible
	if err = mc.handleAuthResult(authData, plugin); err != nil {
		// Authentication failed and MySQL has already closed the connection
		// (https://dev.mysql.com/doc/internals/en/authentication-fails.html).
		// Do not send COM_QUIT, just cleanup and return the error.
		mc.cleanup()
		return nil, err
	}

	if mc.cfg.MaxAllowedPacket > 0 {
		mc.maxAllowedPacket = mc.cfg.MaxAllowedPacket
	} else {
		// Get max allowed packet size
		maxap, err := mc.getSystemVar("max_allowed_packet")
		if err != nil {
			mc.Close()
			return nil, err
		}
		mc.maxAllowedPacket = stringToInt(maxap) - 1
	}
	if mc.maxAllowedPacket < maxPacketSize {
		mc.maxWriteSize = mc.maxAllowedPacket
	}

	// Handle DSN Params
	err = mc.handleParams()
	if err != nil {
		mc.Close()
		return nil, err
	}

	if err = mc.exec("SET @master_binlog_checksum=@@global.binlog_checksum"); err != nil {
		mc.Close()
		return nil, fmt.Errorf("prepareForReplication failed to set @master_binlog_checksum=@@global.binlog_checksum: %v", err)
	}
	return mc, nil
}

// BinLogger implements driver.Pinger interface
func (mc *mysqlConn) ReadBinLog(ctx context.Context, serverID, offset uint32, filename string, flags uint16) (<-chan vtmysql.BinlogEvent, error) {
	fmt.Println(filename, offset, serverID, flags)
	mc.sequence = 0
	length := 4 + //header
		1 + // ComBinlogDump
		4 + // binlog-pos
		2 + // flags
		4 + // server-id
		len(filename) // binlog-filename
	data := make([]byte, length)
	pos := writeByte(data, 4, comBinlogDump)
	pos = writeUint32(data, pos, uint32(offset))
	pos = writeUint16(data, pos, flags)
	pos = writeUint32(data, pos, serverID)
	pos = writeEOFString(data, pos, filename)

	err := mc.writePacket(data)
	if err != nil {
		return nil, err
	}

	var buf []byte
	buf, err = mc.readPacket()
	if err != nil {
		return nil, err
	}

	evCh := make(chan vtmysql.BinlogEvent)
	go func() {
		defer func() {
			close(evCh)
		}()

		for {
			if len(buf) == 0 {
				return
			}
			switch buf[0] {
			case iEOF:
				fmt.Println("iEOF")
				//log.Infof("StartDumpFromBinlogPosition received EOF packet in binlog dump: %#v", buf)
				return
			case iERR:
				fmt.Println("iERR")
				//err := mc.HandleErrorPacket(buf)
				//log.Infof("StartDumpFromBinlogPosition received error packet in binlog dump: %v", err)
				return
			case iOK:
				fmt.Println("OK")
			}
			fmt.Println(buf[1:])
			fmt.Println(vtmysql.NewMysql56BinlogEvent(buf[1:]).NextPosition())

			select {
			case evCh <- vtmysql.NewMysql56BinlogEvent(buf[1:]):
			case <-ctx.Done():
				return
			}
			fmt.Println(vtmysql.NewMysql56BinlogEvent(buf[1:]).NextPosition())

			buf, err = mc.readPacket()
			if err != nil {
				return
			}
		}
	}()

	return evCh, nil
}

func writeEOFString(data []byte, pos int, value string) int {
	pos += copy(data[pos:], value)
	return pos
}

func writeByte(data []byte, pos int, value byte) int {
	data[pos] = value
	return pos + 1
}

func writeUint16(data []byte, pos int, value uint16) int {
	data[pos] = byte(value)
	data[pos+1] = byte(value >> 8)
	return pos + 2
}

func writeUint32(data []byte, pos int, value uint32) int {
	data[pos] = byte(value)
	data[pos+1] = byte(value >> 8)
	data[pos+2] = byte(value >> 16)
	data[pos+3] = byte(value >> 24)
	return pos + 4
}

func writeUint64(data []byte, pos int, value uint64) int {
	data[pos] = byte(value)
	data[pos+1] = byte(value >> 8)
	data[pos+2] = byte(value >> 16)
	data[pos+3] = byte(value >> 24)
	data[pos+4] = byte(value >> 32)
	data[pos+5] = byte(value >> 40)
	data[pos+6] = byte(value >> 48)
	data[pos+7] = byte(value >> 56)
	return pos + 8
}
