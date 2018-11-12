package gosql

import (
	"testing"
)

var multiTableSQLs = []string{
	"select * from (tableA, tableB)",
	"select * from tableA left join tableB on id = id2",
	"select * from tableA right join tableB on id = id2",
	"select * from tableA inner join tableB on id = id2",
	// "select * from tableA full join tableB on tableA.id = tableB.id2",
	"select * from tableA union select * from tableB",
	"select * from tableA union all select * from tableB",
	"select * from tableA union distinct select * from tableB",
}

func TestMultiTableSQLs(t *testing.T) {
	for _, sql := range multiTableSQLs {
		_, err := newQryParser(sql)
		if err == nil {
			t.Error(err)
		}
	}
}

var sqlList = []string{
	// single table
	// "select `id`,ta.name,current_timestamp() from `tableA` as ta",
	// "select ta.`id`, name, group_concat(name2),concat(ta.name3,name4) as b from tableA as ta",
	// "select ta.id, count(1) from tableA as ta",
	// "select ta.ida,tb.idb,tc.idc from `tableA` as ta,`tableB` as tb,`tableC` as tc",
	// "select * from a,(b,c)",
	"select a.*,b.idb from tableA as a inner join tableB as b on a.id = b.id",
	// "select id, ta.name, (select * from tableB) from tableA as ta",
	// "select ta.id from (select * from tableA) as ta",
}

func TestGetSelectors(t *testing.T) {
	for _, sql := range sqlList {
		qp, err := newQryParser(sql)
		if err != nil {
			t.Error(err)
			break
		}
		qp.Parse()
		t.Errorf("%#v", qp)
		qp.print()
	}
}
