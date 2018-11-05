package gomysql

import (
	"fmt"

	vtmysql "vitess.io/vitess/go/mysql"
	"vitess.io/vitess/go/vt/binlog"
	"vitess.io/vitess/go/vt/dbconfigs"
	"vitess.io/vitess/go/vt/vttablet/tabletserver"
	"vitess.io/vitess/go/vt/vttablet/tabletserver/connpool"
	"vitess.io/vitess/go/vt/vttablet/tabletserver/schema"
	"vitess.io/vitess/go/vt/vttablet/tabletserver/tabletenv"
)

const (
	errNewMysqlSlaveConn    = "new mysql slave conn failed, error : %v"
	errNewMysqlSchemaEngine = "Open schema %v engine error : %v"
)

type connCfg struct {
	User     string
	Password string
	IP       string
	Port     int
}

func newConnCfg(user, password, ip string, port int) (cfg *connCfg) {
	cfg = &connCfg{
		IP:       ip,
		Port:     port,
		User:     user,
		Password: password,
	}
	return
}

func (cfg *connCfg) NewSlaveConner() (sc SlaveConner, err error) {
	sc, err = cfg.newSlaveConn()
	return
}

func (cfg *connCfg) newSlaveConn() (sc *slaveConn, err error) {
	sconn, err := binlog.NewSlaveConnection(cfg.formatConnPara())
	if err != nil || sconn == nil {
		err = fmt.Errorf(errNewMysqlSlaveConn, err)
		return
	}

	sc = &slaveConn{
		conn:    sconn,
		schemas: make(map[string]*schema.Engine),
		pbMsgCh: make(chan *binlog.ParsedProto, msgMaxBuffer),
	}
	var se *schemaEngine
	if se, err = cfg.newSchemaEngine(defaultSchema); err != nil {
		return
	}
	var schemaList []string
	if schemaList, err = se.getAllSchemas(); err != nil {
		return
	}
	for _, v := range schemaList {
		if se, err = cfg.newSchemaEngine(v); err != nil {
			return
		}
		sc.schemas[v] = se.engine
	}
	return
}

func (cfg *connCfg) NewSchemaEnginer(schemaName string) (se SchemaEnginer, err error) {
	se, err = cfg.newSchemaEngine(schemaName)
	return
}

func (cfg *connCfg) newSchemaEngine(schemaName string) (se *schemaEngine, err error) {
	cp := cfg.formatConnPara()
	tbServer := tabletserver.NewTabletServerWithNilTopoServer(tabletenv.TabletConfig{})
	engine := schema.NewEngine(tbServer, tabletenv.TabletConfig{})
	dbcfgs := dbconfigs.NewTestDBConfigs(*cp, *cp, schemaName)
	engine.InitDBConfig(dbcfgs)
	err = engine.Open()
	if err != nil {
		err = fmt.Errorf(errNewMysqlSchemaEngine, schemaName, err)
		return nil, err
	}
	// init pool part
	pools := connpool.New("", 3, 0, tbServer)
	dbaParams := dbcfgs.DbaWithDB()
	pools.Open(dbaParams, dbaParams, dbaParams)
	se = &schemaEngine{
		engine: engine,
		pools:  pools,
	}
	return se, nil
}

func (cfg *connCfg) formatConnPara() *vtmysql.ConnParams {
	return &vtmysql.ConnParams{
		Host:  cfg.IP,
		Port:  cfg.Port,
		Uname: cfg.User,
		Pass:  cfg.Password,
	}
}
