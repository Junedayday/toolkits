package gomysql

import "fmt"

var (
	errNewMysqlSlaveConn      = "new msyql slave conn failed, error : %v"
	errNewMysqlSchemaEngine   = "Open schema %v engine error : %v"
	errNewMysqlSchemaNotExist = "schema %v may not exist error : %v"

	errConnClosed = fmt.Errorf("connection is closed")
)
