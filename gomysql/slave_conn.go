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

	errConnClosed  = "connection is closed"
	errMsgChClosed = "message channel is closed"
)

type slaveConn struct {
	conn    *binlog.SlaveConnection
	schemas map[string]*schema.Engine
	pbMsgCh chan *binlog.ParsedProto
	pos     vtmysql.Position
	sync.RWMutex
}

func (sc *slaveConn) CurrentDumpBinlogAsPB(dealFunc func(*pbmysql.Event), reloadCh chan struct{}) (err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var eventCh <-chan vtmysql.BinlogEvent
	_, eventCh, err = sc.conn.StartBinlogDumpFromCurrent(ctx)
	if err != nil {
		glog.Fatalf("StartBinlogDumpFromBinlogBeforeTimestamp %v", err)
	}
	go sc.startProtoParsing(ctx, eventCh)
	return sc.dealProtoMsg(dealFunc, reloadCh)
}

func (sc *slaveConn) ContinueDumpBinlogAsPB(dealFunc func(*pbmysql.Event), reloadCh chan struct{}) (err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var eventCh <-chan vtmysql.BinlogEvent
	eventCh, err = sc.conn.StartBinlogDumpFromPosition(ctx, sc.pos)
	if err != nil {
		err = fmt.Errorf("StartBinlogDumpFromPosition %v", err)
		return
	}
	go sc.startProtoParsing(ctx, eventCh)
	return sc.dealProtoMsg(dealFunc, reloadCh)
}

func (sc *slaveConn) GetPosition() (b []byte) {
	b = []byte(vtmysql.EncodePosition(sc.pos))
	return
}

func (sc *slaveConn) SetPosition(b []byte) (err error) {
	sc.pos, err = vtmysql.DecodePosition(string(b))
	return
}

func (sc *slaveConn) Close() {
	close(sc.pbMsgCh)
	sc.conn.Close()
}

func (sc *slaveConn) startProtoParsing(ctx context.Context, eventCh <-chan vtmysql.BinlogEvent) {
	// ParseEvents is packed in "vitess.io/vitess/go/vt/binlog"
	binlog.ParseProtoEvents(ctx, eventCh, sc.schemas, sc.pbMsgCh, sc.pos, true)
}

func (sc *slaveConn) dealProtoMsg(dealFunc func(*pbmysql.Event), reloadCh chan struct{}) (err error) {
	for {
		select {
		case msg, ok := <-sc.pbMsgCh:
			if ok {
				if msg.ErrInfo != nil {
					if msg.ErrInfo == binlog.ErrServerEOF {
						return fmt.Errorf(errConnClosed)
					}
					glog.Errorf("parse event error : %v", msg.ErrInfo)
				} else if len(msg.ReloadDDL) != 0 {
					glog.Warningf("Find a DDL %v, going to reload", msg.ReloadDDL)
					reloadCh <- struct{}{}
					return
				} else {
					for _, v := range msg.Events {
						glog.Infof("Get a binlog event %v-%v-%v, type %v\n", v.Schema, v.Table, v.Id, v.Et)
						dealFunc(v)
					}
					sc.savePosition(msg.NextPos)
				}
			} else {
				glog.Warning("message channel is closed in slave connection")
				return fmt.Errorf(errMsgChClosed)
			}
		}
	}
}

func (sc *slaveConn) savePosition(position vtmysql.Position) {
	sc.Lock()
	sc.pos = position
	sc.Unlock()
}
