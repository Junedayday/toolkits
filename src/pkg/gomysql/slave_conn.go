package gomysql

import (
	"context"
	"fmt"
	"proto/mysql/pbmysql"
	"sync"

	"github.com/golang/glog"
	vtmysql "vitess.io/vitess/go/mysql"
	"vitess.io/vitess/go/vt/binlog"
	"vitess.io/vitess/go/vt/vttablet/tabletserver/schema"
)

const (
	msgMaxBuffer = 100
)

type slaveConn struct {
	conn    *binlog.SlaveConnection
	schemas map[string]*schema.Engine
	// binlog event transfer to sql statement msg
	statMsgCh chan *binlog.ParsedStat
	// binlog event transfer to proto defined struct message
	pbMsgCh chan *binlog.ParsedProto
	// record postion
	pos vtmysql.Position
	sync.RWMutex
}

func newSlaveConn(cfg *connCfg) (sc *slaveConn, err error) {
	cp := cfg.formatCP()
	sconn, err := binlog.NewSlaveConnection(cp)
	if err != nil || sconn == nil {
		err = fmt.Errorf(errNewMysqlSlaveConn, err)
		return
	}
	sc = &slaveConn{
		conn:      sconn,
		schemas:   make(map[string]*schema.Engine),
		statMsgCh: make(chan *binlog.ParsedStat, msgMaxBuffer),
		pbMsgCh:   make(chan *binlog.ParsedProto, msgMaxBuffer),
	}

	// get all schemas
	var engine *schemaEngine
	if engine, err = newSchemaEngine(cfg, defaultSchema); err != nil {
		return
	}
	var schemas []string
	if schemas, err = engine.GetAllSchemas(); err != nil {
		return
	}
	for _, v := range schemas {
		if engine, err = newSchemaEngine(cfg, v); err != nil {
			return
		}
		sc.schemas[v] = engine.engine
	}
	return
}

func (sc *slaveConn) StartBinlogDumpFromCurrentAsProto(dealFunc func(*pbmysql.Event), reloadCh chan struct{}) (err error) {
	ctx, cancel := context.WithCancel(context.Background())
	var eventCh <-chan vtmysql.BinlogEvent
	var pos vtmysql.Position
	pos, eventCh, err = sc.conn.StartBinlogDumpFromCurrent(ctx)
	// save position
	sc.savePosition(pos)
	if err != nil {
		glog.Fatalf("StartBinlogDumpFromBinlogBeforeTimestamp %v", err)
	}
	go sc.startProtoParsing(ctx, eventCh)
	defer cancel()
	return sc.dealProtoMsg(dealFunc, reloadCh)
}

func (sc *slaveConn) StartBinlogDumpFromPositionAsProto(dealFunc func(*pbmysql.Event), reloadCh chan struct{}) (err error) {
	ctx, cancel := context.WithCancel(context.Background())
	var eventCh <-chan vtmysql.BinlogEvent
	eventCh, err = sc.conn.StartBinlogDumpFromPosition(ctx, sc.pos)
	if err != nil {
		err = fmt.Errorf("StartBinlogDumpFromPosition %v", err)
		return
	}
	go sc.startProtoParsing(ctx, eventCh)
	defer cancel()
	return sc.dealProtoMsg(dealFunc, reloadCh)
}

func (sc *slaveConn) StartBinlogDumpFromCurrentAsStat(dealFunc func(binlog.FullBinlogStatement)) (err error) {
	ctx := context.Background()
	var eventCh <-chan vtmysql.BinlogEvent
	var pos vtmysql.Position
	pos, eventCh, err = sc.conn.StartBinlogDumpFromCurrent(ctx)
	// save position
	sc.savePosition(pos)
	if err != nil {
		err = fmt.Errorf("StartBinlogDumpFromBinlogBeforeTimestamp %v", err)
		return
	}
	go sc.startStatsParsing(ctx, eventCh)

	return sc.dealStatsMsg(dealFunc)
}

func (sc *slaveConn) StartBinlogDumpFromPositionAsStat(dealFunc func(binlog.FullBinlogStatement)) (err error) {
	ctx := context.Background()
	var eventCh <-chan vtmysql.BinlogEvent
	eventCh, err = sc.conn.StartBinlogDumpFromPosition(ctx, sc.pos)
	if err != nil {
		err = fmt.Errorf("StartBinlogDumpFromPosition %v", err)
		return
	}
	go sc.startStatsParsing(ctx, eventCh)

	return sc.dealStatsMsg(dealFunc)
}

func (sc *slaveConn) EncodePosition() (b []byte) {
	b = []byte(vtmysql.EncodePosition(sc.pos))
	return
}

func (sc *slaveConn) DecodePosition(b []byte) (err error) {
	sc.pos, err = vtmysql.DecodePosition(string(b))
	return
}

func (sc *slaveConn) Close() {
	sc.conn.Close()
}

func (sc *slaveConn) startProtoParsing(ctx context.Context, eventCh <-chan vtmysql.BinlogEvent) {
	// ParseEvents is packed in "vitess.io/vitess/go/vt/binlog"
	binlog.ParseProtoEvents(ctx, eventCh, sc.schemas, sc.pbMsgCh, sc.pos, true)
}

func (sc *slaveConn) startStatsParsing(ctx context.Context, eventCh <-chan vtmysql.BinlogEvent) {
	// ParseEvents is packed in "vitess.io/vitess/go/vt/binlog"
	binlog.ParseStatsEvents(ctx, eventCh, sc.schemas, sc.statMsgCh)
}

func (sc *slaveConn) dealProtoMsg(dealFunc func(*pbmysql.Event), reloadCh chan struct{}) (err error) {
	for {
		select {
		case msg, ok := <-sc.pbMsgCh:
			if ok {
				if msg.ErrInfo != nil {
					if msg.ErrInfo == binlog.ErrServerEOF {
						return errConnClosed
					}
					glog.Errorf("parse event error : %v", msg.ErrInfo)
				} else if msg.Reload {
					glog.Infof("Find a DDL, going to reload table")
					reloadCh <- struct{}{}
					fmt.Println("send reload")
					return
				} else {
					for _, v := range msg.Events {
						dealFunc(v)
					}
					sc.savePosition(msg.NextPos)
				}
			}
		}
	}
}

func (sc *slaveConn) dealStatsMsg(dealFunc func(binlog.FullBinlogStatement)) (err error) {
	for {
		select {
		case msg, ok := <-sc.statMsgCh:
			if ok {
				if msg.ErrInfo != nil {
					if msg.ErrInfo == binlog.ErrServerEOF {
						return errConnClosed
					}
					glog.Errorf("parse event error : %v", msg.ErrInfo)
				} else {
					for _, v := range msg.Stats {
						dealFunc(v)
					}
					sc.savePosition(msg.NextPos)
				}
			}
		}
	}
}

func (sc *slaveConn) SetMasterPosition() (err error) {
	sc.pos, err = sc.conn.Conn.MasterPosition()
	return
}

func (sc *slaveConn) savePosition(position vtmysql.Position) {
	sc.Lock()
	sc.pos = position
	sc.Unlock()
}
