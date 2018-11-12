package gosql

import (
	"fmt"

	"vitess.io/vitess/go/vt/sqlparser"
)

type tableInner struct {
	sourceTable string
	sourceCol   string
	aliasCol    string
}

type tablerSet struct {
	// used in subquery with a alias
	// checked src
	checkedSrc []string
	// aliasSet   string
	// in         []tableInner
	// used in sqlparser.TableName type
	src      string
	aliasSrc string
}

func newTablerWithName(tableName, alias string) tablerSet {
	return tablerSet{
		src:      tableName,
		aliasSrc: alias,
	}
}

func parseTablersExpr(node sqlparser.TableExpr) (tbers []tablerSet) {
	switch subNode := node.(type) {
	case *sqlparser.AliasedTableExpr:
		tbers = append(tbers, parseAliasedTableExpr(subNode)...)
	case *sqlparser.ParenTableExpr:
		tbers = append(tbers, parseParenTableExpr(subNode)...)
	case *sqlparser.JoinTableExpr:
		tbers = append(tbers, parseJoinTableExpr(subNode)...)
	default:
		fmt.Printf("%#v", subNode)
		tbers = append(tbers, tablerSet{src: "unsupport"})
	}
	return
}

func parseParenTableExpr(node *sqlparser.ParenTableExpr) (tbers []tablerSet) {
	for _, v := range node.Exprs {
		tbers = append(tbers, parseTablersExpr(v)...)
	}
	return
}

func parseJoinTableExpr(node *sqlparser.JoinTableExpr) (tbers []tablerSet) {
	// tber := newTablerWithName(subNode.Name.String(), node.As.String())
	tbers = append(tbers, parseTablersExpr(node.LeftExpr)...)
	tbers = append(tbers, parseTablersExpr(node.RightExpr)...)
	return
}

func parseAliasedTableExpr(node *sqlparser.AliasedTableExpr) (tbers []tablerSet) {
	switch subNode := node.Expr.(type) {
	case sqlparser.TableName:
		tber := newTablerWithName(subNode.Name.String(), node.As.String())
		tbers = append(tbers, tber)
	case *sqlparser.Subquery:
		panic("not support for sub query")
	}
	return
}

// func (tber *tabler) addTableCol(tableName, col string) {
// 	var isStar bool
// 	if col == constStarSelect {
// 		isStar = true
// 	}

// 	for k, v := range tber.in {
// 		if v.tableName == tableName {
// 			if isStar {
// 				tber.in[k].isStar = true
// 				tber.in[k].cols = []string{}
// 			} else {
// 				tber.in[k].cols = append(tber.in[k].cols, col)
// 			}
// 			return
// 		}
// 	}

// 	if isStar {
// 		tber.in = append(tber.in, tableInner{tableName: tableName, isStar: true})
// 	} else {
// 		tber.in = append(tber.in, tableInner{tableName: tableName, cols: []string{col}})
// 	}
// }
