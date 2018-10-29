package gomysql

import (
	"proto/mysql/pbmysql"

	"vitess.io/vitess/go/vt/binlog"
)

// ConnCfger implement for msyql instance conn
type ConnCfger interface {
	NewSlaveConner() (sc SlaveConner, err error)
	NewSchemaEnginer(schemaName string) (se SchemaEnginer, err error)
}

// NewConnCfger new a configuration to mysql instance
func NewConnCfger(user, password, ip string, port int) ConnCfger {
	return newConnCfg(user, password, ip, port)
}

// SlaveConner implement for slave mysql conn
type SlaveConner interface {
	StartBinlogDumpFromCurrentAsProto(cdealFunc func(*pbmysql.Event), reloadCh chan struct{}) (err error)
	StartBinlogDumpFromPositionAsProto(dealFunc func(*pbmysql.Event), reloadCh chan struct{}) (err error)
	StartBinlogDumpFromCurrentAsStat(dealFunc func(binlog.FullBinlogStatement)) (err error)
	StartBinlogDumpFromPositionAsStat(dealFunc func(binlog.FullBinlogStatement)) (err error)
	EncodePosition() (b []byte)
	DecodePosition(b []byte) (err error)
	SetMasterPosition() (err error)
	Close()
}

// SchemaEnginer implement for connection to a schema in mysql
type SchemaEnginer interface {
	GetAllSchemas() (schemaNames []string, err error)
	Close()
}
