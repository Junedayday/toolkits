package gomysql

import (
	"fmt"

	"vitess.io/vitess/go/sqltypes"
	"vitess.io/vitess/go/vt/vttablet/tabletserver/connpool"

	"vitess.io/vitess/go/vt/dbconfigs"
	"vitess.io/vitess/go/vt/vttablet/tabletserver"
	"vitess.io/vitess/go/vt/vttablet/tabletserver/schema"
	"vitess.io/vitess/go/vt/vttablet/tabletserver/tabletenv"
)

type schemaEngine struct {
	engine *schema.Engine
	pools  *connpool.Pool
}

func newSchemaEngine(cfg *connCfg, schemaName string) (se *schemaEngine, err error) {
	// init engine part
	cp := cfg.formatCP()
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

func (se *schemaEngine) Close() {
	se.engine.Close()
}

func (se *schemaEngine) GetAllSchemas() (schemaNames []string, err error) {
	var tableData *sqltypes.Result
	tableData, err = se.getTableData(showAllDatabases, defaultMaxRows, false)
	for _, row := range tableData.Rows {
		schemaNames = append(schemaNames, row[0].ToString())
	}
	return
}

// used to get query result in database
func (se *schemaEngine) getTableData(query string, maxrows int, wantfields bool) (*sqltypes.Result, error) {
	ctx := tabletenv.LocalContext()
	conn, err := se.pools.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Recycle()

	return conn.Exec(ctx, query, maxrows, wantfields)
}
