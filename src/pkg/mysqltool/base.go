package mysqltool

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Mysqlcfg Default cfg
type Mysqlcfg struct {
	IP       string
	Port     int
	User     string
	Password string
}

// SWMysql use for mysql operations
type swMysql struct {
	*sql.DB
	connStr string
}

func newMysqlPool(ctx context.Context, cfg *Mysqlcfg, dbname string) (SWMysql, error) {
	if dbname == "" {
		dbname = "information_schema"
	}
	swdb := &swMysql{
		connStr: fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8", cfg.User, cfg.Password, cfg.IP, cfg.Port, dbname),
	}
	db, err := sql.Open("mysql", swdb.connStr)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)
	swdb.DB = db
	// reconnect to mysql
	go checkConn(ctx, swdb)
	return swdb, err
}

func checkConn(ctx context.Context, swdb *swMysql) {
	for {
		time.Sleep(time.Second)
		select {
		case <-ctx.Done():
			swdb.Close()
		default:
			if swdb.Ping() != nil {
				reConn(swdb)
			}
			time.Sleep(4 * time.Second)
		}
	}
}

func reConn(swdb *swMysql) {
	db, err := sql.Open("mysql", swdb.connStr)
	if err != nil {
		return
	}
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)
	swdb.DB = db
}

// queryResult is used for packing sql results
type queryResult struct {
	Data    [][]string // query results
	Columns []string   // column nums of result
	Num     int        // sql returns nums
}

func (result *queryResult) HasData() bool {
	return len(result.Data) > 0 && len(result.Columns) > 0 && result.Num > 0
}

func (result *queryResult) GetData() [][]string {
	return result.Data
}

func (result *queryResult) GetColumns() []string {
	return result.Columns
}

func (result *queryResult) GetNum() int {
	return result.Num
}

func initQueryResult() *queryResult {
	return &queryResult{
		Data:    make([][]string, 0, 0),
		Columns: make([]string, 0),
	}
}
