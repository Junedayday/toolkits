package mysqltool

import (
	"database/sql"
)

func (db *swMysql) QueryToResult(sqls string, args ...interface{}) (QueryResult, error) {
	rows, err := db.Query(sqls, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// init query result
	result := initQueryResult()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	result.Columns, err = rows.Columns()
	if err != nil {
		return nil, err
	}

	newVal := make([]sql.RawBytes, 0)
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		newVal = append(newVal, values...)
		var tmp []string
		for _, col := range values {
			tmp = append(tmp, string(col))
		}
		if err != nil {
			return nil, err
		}
		result.Data = append(result.Data, tmp)
		result.Num++
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (db *swMysql) ISGetSchemaTableLists() (QueryResult, error) {
	return db.QueryToResult(constQuerySchemaTables)
}
