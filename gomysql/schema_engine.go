package gomysql

import (
	"vitess.io/vitess/go/sqltypes"
	"vitess.io/vitess/go/vt/vttablet/tabletserver/connpool"

	"vitess.io/vitess/go/vt/vttablet/tabletserver/schema"
	"vitess.io/vitess/go/vt/vttablet/tabletserver/tabletenv"
)

type schemaEngine struct {
	engine *schema.Engine
	pools  *connpool.Pool
}

func (se *schemaEngine) Close() {
	se.engine.Close()
}

func (se *schemaEngine) getAllSchemas() (schemaNames []string, err error) {
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
