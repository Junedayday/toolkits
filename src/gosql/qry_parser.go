package gosql

import (
	"fmt"

	"vitess.io/vitess/go/vt/sqlparser"
)

const constStarSelect = "*"

type qryParser struct {
	sql string
	// use type in sqlparser
	// Union,Select,Stream,Insert,Update,Delete,Set
	// DBDDL,DDL,Show,Use,Begin,Commit,Rollback,OtherRead,OtherAdmin
	stmt      sqlparser.Statement
	tablers   []tablerSet
	selectors []selector
}

func (qp *qryParser) Parse() {
	qp.selectors = qp.GetSelectors()
	qp.tablers = qp.GetTablers()
}

func (qp *qryParser) print() {
	for _, v := range qp.selectors {
		if v.from == "" {
			if len(qp.tablers) == 1 && qp.tablers[0].src != "" {
				fmt.Println("col :", v.col, ",table :", qp.tablers[0].src)
			} else {
				fmt.Println("col :", v.col, ",table :", "unknown")
			}
		} else {
			var src string
			for _, v2 := range qp.tablers {
				if v.from == v2.aliasSrc || (v2.aliasSrc == "" && v.from == v2.src) {
					src = v2.src
					break
				}
			}
			fmt.Println("col :", v.col, ",table :", src)
		}
	}
}

func (qp *qryParser) GetSelectors() (selectors []selector) {
	selectAST, _ := qp.stmt.(*sqlparser.Select)
	for _, node := range selectAST.SelectExprs {
		selectors = append(selectors, parseSelectExpr(node)...)
	}
	return
}

func (qp *qryParser) GetTablers() (tablers []tablerSet) {
	selectAST, _ := qp.stmt.(*sqlparser.Select)
	for _, node := range selectAST.From {
		tablers = append(tablers, parseTablersExpr(node)...)
	}
	return
}

func newQryParser(originSQL string) (qp *qryParser, err error) {
	var stmt sqlparser.Statement
	stmt, err = getStmt(originSQL)
	if err != nil {
		return
	}
	// if isMultiTable(stmt) {
	// 	err = fmt.Errorf("not support for multitable %v", originSQL)
	// 	return
	// }
	return &qryParser{
		sql:  originSQL,
		stmt: stmt,
	}, nil
}

func getStmt(sql string) (sqlparser.Statement, error) {
	if sqlparser.Preview(sql) != sqlparser.StmtSelect {
		return nil, fmt.Errorf("%v not a select", sql)
	}
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		return nil, fmt.Errorf("%v parse sql failed %v", sql, err)
	}
	return stmt, nil
}

// mutitable not support temporary
// full join not support

// func isMultiTable(stmt sqlparser.Statement) bool {
// 	_, ok := stmt.(*sqlparser.Union)
// 	if ok {
// 		return true
// 	}

// 	astSelect, ok := stmt.(*sqlparser.Select)
// 	if !ok {
// 		return true
// 	}

// 	for _, v := range astSelect.From {
// 		if _, ok = v.(*sqlparser.ParenTableExpr); ok {
// 			return true
// 		}
// 		if _, ok = v.(*sqlparser.JoinTableExpr); ok {
// 			return true
// 		}
// 	}
// 	return false
// }
