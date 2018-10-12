package gomysql

import (
	"pkg/tcfgs"
	"proto/mysql/pbmysql"
	"testing"
	"time"

	"vitess.io/vitess/go/vt/binlog"
)

type MysqlCfg struct {
	IP       string
	Port     int
	User     string
	Password string
}

type yamlCfg struct {
	Mysql MysqlCfg
}

func testNewConnCfg() (ConnCfger, error) {
	mysqlCfg, err := tcfgs.GetTestMysqlCfg()
	if err != nil {
		return nil, err
	}
	return NewConnCfger(mysqlCfg.User, mysqlCfg.Password, mysqlCfg.IP, mysqlCfg.Port), nil
}

func TestMysqlSlaveConn(t *testing.T) {
	cfg, err := testNewConnCfg()
	if err != nil {
		t.Errorf("new config error %v", err)
		return
	}
	sc, err := cfg.NewSlaveConner()
	if err != nil {
		t.Errorf("new slave connection error %v", err)
		return
	}
	sc.Close()

	cfg = NewConnCfger("root", "123123", "127.0.0.1", 3306)
	_, err = cfg.NewSlaveConner()
	if err == nil {
		t.Errorf("new slave connection to unkown host")
		return
	}
}

func TestMysqlSchemaEngine(t *testing.T) {
	cfg, err := testNewConnCfg()
	if err != nil {
		t.Errorf("new config error %v", err)
		return
	}
	se, err := cfg.NewSchemaEnginer(defaultSchema)
	if err != nil {
		t.Errorf("new schema engine error %v", err)
		return
	}
	se.Close()

	cfg2, err := testNewConnCfg()
	if err != nil {
		t.Errorf("new config error %v", err)
		return
	}
	_, err2 := cfg2.NewSchemaEnginer("not_exit")
	if err2 == nil {
		t.Errorf("new not_exit schema engine error")
		return
	}
}

func TestMysqlGetAllSchemas(t *testing.T) {
	cfg, err := testNewConnCfg()
	if err != nil {
		t.Errorf("new config error %v", err)
		return
	}
	se, err := cfg.NewSchemaEnginer(defaultSchema)
	if err != nil {
		t.Errorf("new schema engine error %v", err)
		return
	}
	schemas, err := se.GetAllSchemas()
	if err != nil || len(schemas) == 0 {
		t.Errorf("get all schema engine error %v", err)
		return
	}
	se.Close()
}

func TestMysqlGetBinLogStats(t *testing.T) {
	cfg, err := testNewConnCfg()
	if err != nil {
		t.Errorf("new config error %v", err)
		return
	}
	sc, err := cfg.NewSlaveConner()
	if err != nil {
		t.Errorf("new slave connection error %v", err)
		return
	}
	defer sc.Close()

	f := func(msg binlog.FullBinlogStatement) {}
	go sc.StartBinlogDumpFromCurrentAsStat(f)
	time.Sleep(time.Second)
}

func TestMysqlGetBinLogProto(t *testing.T) {
	cfg, err := testNewConnCfg()
	if err != nil {
		t.Errorf("new config error %v", err)
		return
	}
	sc, err := cfg.NewSlaveConner()
	if err != nil {
		t.Errorf("new slave connection error %v", err)
		return
	}
	defer sc.Close()

	f := func(msg *pbmysql.Event) {}
	go sc.StartBinlogDumpFromCurrentAsProto(f)
	time.Sleep(time.Second)
}

func TestMysqlGetBinLogStatsFromPos(t *testing.T) {
	cfg, err := testNewConnCfg()
	if err != nil {
		t.Errorf("new config error %v", err)
		return
	}
	sc, err := cfg.NewSlaveConner()
	if err != nil {
		t.Errorf("new slave connection error %v", err)
		return
	}
	defer sc.Close()

	f := func(msg binlog.FullBinlogStatement) {}
	err = sc.SetMasterPosition()
	if err != nil {
		t.Errorf("set position error %v", err)
		return
	}
	go sc.StartBinlogDumpFromPositionAsStat(f)
	time.Sleep(time.Second)
}

func TestMysqlGetBinLogProtoFromPos(t *testing.T) {
	cfg, err := testNewConnCfg()
	if err != nil {
		t.Errorf("new config error %v", err)
		return
	}
	sc, err := cfg.NewSlaveConner()
	if err != nil {
		t.Errorf("new slave connection error %v", err)
		return
	}
	defer sc.Close()

	f := func(msg *pbmysql.Event) {}
	err = sc.SetMasterPosition()
	if err != nil {
		t.Errorf("set position error %v", err)
		return
	}
	go sc.StartBinlogDumpFromPositionAsProto(f)
	time.Sleep(time.Second)
}
