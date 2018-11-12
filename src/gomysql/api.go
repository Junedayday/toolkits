package gomysql

import (
	"proto/mysql/pbmysql"

	"vitess.io/vitess/go/sqltypes"
)

// ConnCfger implement for a msyql instance conn
type ConnCfger interface {
	NewSlaveConner() (sc SlaveConner, err error)
	NewSchemaEnginer(schemaName string) (se SchemaEnginer, err error)
}

// NewConnCfger : user,password,ip,port are essential info
func NewConnCfger(user, password, ip string, port int) ConnCfger {
	return newConnCfg(user, password, ip, port)
}

// SlaveConner implement for slave mysql conn
type SlaveConner interface {
	CurrentDumpBinlogAsPB(dealFunc func(*pbmysql.Event), reloadCh chan struct{}) (err error)
	ContinueDumpBinlogAsPB(dealFunc func(*pbmysql.Event), reloadCh chan struct{}) (err error)
	GetPosition() (b []byte)
	SetPosition(b []byte) (err error)
	Close()
}

// SchemaEnginer implement for connection to a schema in mysql
type SchemaEnginer interface {
	QueryToStruct(query string, maxrows int, wantfields bool) (*sqltypes.Result, error)
	Close()
}
