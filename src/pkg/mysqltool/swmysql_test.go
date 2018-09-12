package mysqltool

import (
	"context"
	"pkg/cfgtool"
	"testing"
)

type dbInfo struct {
	IP       string
	Port     int
	User     string
	Password string
}

type dbCfg struct {
	Mysql1 *Mysqlcfg
	Mysql2 *Mysqlcfg
}

func TestMysqlConn(t *testing.T) {
	cfg := &dbCfg{
		Mysql1: &Mysqlcfg{},
		Mysql2: &Mysqlcfg{},
	}
	err := cfgtool.LoadCfgFromYamlFile("../../configs/testing/db.yaml", cfg)
	if err != nil {
		t.Error("read yaml failed")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err = NewMysqlPool(ctx, cfg.Mysql1, "")
	if err != nil {
		t.Errorf("connect to mysql1 failed %v", err)
		return
	}

	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()
	_, err = NewMysqlPool(ctx2, cfg.Mysql2, "")
	if err != nil {
		t.Errorf("connect to mysql2 failed %v", err)
		return
	}
}

func TestGetSchemaTableListSlice(t *testing.T) {
	cfg := &dbCfg{
		Mysql1: &Mysqlcfg{},
		Mysql2: &Mysqlcfg{},
	}
	err := cfgtool.LoadCfgFromYamlFile("../../configs/testing/db.yaml", cfg)
	if err != nil {
		t.Error("read yaml failed")
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := newMysqlPool(ctx, cfg.Mysql1, "")
	if err != nil {
		t.Errorf("connect to mysql1 failed %v", err)
		return
	}
	result, err := db.ISGetSchemaTableLists()
	if err != nil {
		t.Errorf("query error %v", err)
	} else if !result.HasData() {
		t.Error("result has no data")
	}
}
