package gomysql

import (
	vtmysql "vitess.io/vitess/go/mysql"
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
	sc, err = newSlaveConn(cfg)
	return
}

func (cfg *connCfg) NewSchemaEnginer(schemaName string) (se SchemaEnginer, err error) {
	se, err = newSchemaEngine(cfg, schemaName)
	return
}

// func (cfg *connCfg) formatDSN(schemaName string) (dsn string) {
// 	dsn = fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", cfg.User, cfg.Password, cfg.IP, cfg.Port, schemaName)
// 	return
// }

func (cfg *connCfg) formatCP() *vtmysql.ConnParams {
	return &vtmysql.ConnParams{
		Host:  cfg.IP,
		Port:  cfg.Port,
		Uname: cfg.User,
		Pass:  cfg.Password,
	}
}
