package mysqltool

import (
	"context"

	_ "github.com/go-sql-driver/mysql" //import mysql driver
)

const (
	// IS : db information schema
	IS = "information_schema"
)

// QueryResult is interface of queryResult
type QueryResult interface {
	// check if the result has data
	HasData() bool
	// get data
	GetData() [][]string
	GetColumns() []string
	GetNum() int
}

// SWMysql is interface of swmsyql
type SWMysql interface {
	// common use
	QueryToResult(sql string, args ...interface{}) (QueryResult, error)

	// query certain dbs:
	// information_schema : IS
	ISGetSchemaTableLists() (QueryResult, error)
}

// NewMysqlPool : init mysql, default db is sys
func NewMysqlPool(ctx context.Context, cfg *Mysqlcfg, dbname string) (SWMysql, error) {
	return newMysqlPool(ctx, cfg, dbname)
}
