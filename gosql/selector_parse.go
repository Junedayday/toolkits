package gosql

import (
	"fmt"
	"strings"

	"vitess.io/vitess/go/vt/sqlparser"
)

// format like : `from`.`col` (as) `alias`
type selector struct {
	col  string
	from string
}

func newSelector(col, from, alias string) selector {
	return selector{
		col:  col,
		from: from,
	}
}

func parseSelectExpr(node sqlparser.SelectExpr) (selectors []selector) {
	switch node := node.(type) {
	case *sqlparser.AliasedExpr:
		selectors = append(selectors, parseAliasedExpr(node)...)
	case *sqlparser.StarExpr:
		selectors = append(selectors, parseStarExpr(node)...)
	}
	return
}

func parseStarExpr(node *sqlparser.StarExpr) (selectors []selector) {
	slctor := newSelector(constStarSelect, node.TableName.Name.String(), node.TableName.Qualifier.String())
	selectors = append(selectors, slctor)
	return
}

func parseAliasedExpr(node *sqlparser.AliasedExpr) (selectors []selector) {
	switch subNode := node.Expr.(type) {
	case *sqlparser.ColName:
		selectors = append(selectors, parseColName(subNode)...)
	case *sqlparser.FuncExpr:
		selectors = append(selectors, parseFuncExpr(subNode)...)
	case *sqlparser.SQLVal:
		// ignore for const sql value
		parseSQLVal(subNode)
	case *sqlparser.GroupConcatExpr:
		selectors = append(selectors, parseGroupConcatExpr(subNode)...)
	case *sqlparser.Subquery:
		// not support now!
		selectors = append(selectors, parseSubquery(subNode)...)
	default:
		fmt.Printf("%#v", subNode)
		selectors = append(selectors, selector{col: "unsupport"})
	}
	return
}

func parseSQLVal(node *sqlparser.SQLVal) (selectors []selector) {
	return
}

func parseColName(node *sqlparser.ColName) (selectors []selector) {
	col := strings.ToLower(node.Name.String())
	from := strings.ToLower(node.Qualifier.Name.String())
	selectors = append(selectors, selector{col: col, from: from})
	return
}

func parseFuncExpr(node *sqlparser.FuncExpr) (selectors []selector) {
	for _, v := range node.Exprs {
		selectors = append(selectors, parseSelectExpr(v)...)
	}
	return
}

func parseGroupConcatExpr(node *sqlparser.GroupConcatExpr) (selectors []selector) {
	for _, v := range node.Exprs {
		selectors = append(selectors, parseSelectExpr(v)...)
	}
	return
}

func parseSubquery(node *sqlparser.Subquery) (selectors []selector) {
	panic("not support for sub query now !")
}
