package golangMysqlPool

import (
	"fmt"
	"log"
	"strings"
)

var (
	supportFunction = []string{
		"NOW()", "CURRENT_TIMESTAMP", "UUID()", "RAND()", "CURDATE()",
		"CURTIME()", "UNIX_TIMESTAMP()", "UTC_TIMESTAMP()", "SYSDATE()",
		"LOCALTIME()", "LOCALTIMESTAMP()", "PI()", "DATABASE()", "USER()",
		"VERSION()",
	}
)

func (db *Pool) DB(databaseName string) *QueryBuilder {
	_, err := db.db.Exec(fmt.Sprintf("USE `%s`", databaseName))
	if err != nil {
		db.Logger.Action(true,
			fmt.Sprintf("Failed to switch to database %s", databaseName),
			err.Error(),
		)
	}

	return &QueryBuilder{
		db:         db.db,
		Database:   &databaseName,
		SelectList: []string{"*"},
		Logger:     db.Logger,
	}
}

func (q *QueryBuilder) Table(tableName string) *QueryBuilder {
	q.TableName = &tableName
	return q
}

func (q *QueryBuilder) Select(fields ...string) *QueryBuilder {
	if len(fields) > 0 {
		q.SelectList = fields
	}
	return q
}

func (q *QueryBuilder) Total() *QueryBuilder {
	q.WithTotal = true
	return q
}

func (q *QueryBuilder) InnerJoin(table, first, operator string, second ...string) *QueryBuilder {
	return q.join("INNER", table, first, operator, second...)
}

func (q *QueryBuilder) LeftJoin(table, first, operator string, second ...string) *QueryBuilder {
	return q.join("LEFT", table, first, operator, second...)
}

func (q *QueryBuilder) RightJoin(table, first, operator string, second ...string) *QueryBuilder {
	return q.join("RIGHT", table, first, operator, second...)
}

// * private method
func (q *QueryBuilder) join(joinType, table, first, operator string, second ...string) *QueryBuilder {
	var secondField string
	if len(second) > 0 {
		secondField = second[0]
	} else {
		secondField = operator
		operator = "="
	}

	if !strings.Contains(first, ".") {
		first = fmt.Sprintf("`%s`", first)
	}
	if !strings.Contains(secondField, ".") {
		secondField = fmt.Sprintf("`%s`", secondField)
	}

	joinClause := fmt.Sprintf("%s JOIN `%s` ON %s %s %s", joinType, table, first, operator, secondField)
	q.JoinList = append(q.JoinList, joinClause)
	return q
}

func (q *QueryBuilder) Where(column string, operator interface{}, value ...interface{}) *QueryBuilder {
	var targetValue interface{}
	var targetOperator string

	if len(value) == 0 {
		targetValue = operator
		targetOperator = "="
	} else {
		targetOperator = fmt.Sprintf("%v", operator)
		targetValue = value[0]
	}

	if targetOperator == "LIKE" {
		if str, ok := targetValue.(string); ok {
			targetValue = fmt.Sprintf("%%%s%%", str)
		}
	}

	if !strings.Contains(column, "(") && !strings.Contains(column, ".") {
		column = fmt.Sprintf("`%s`", column)
	}

	placeholder := "?"
	if targetOperator == "IN" {
		placeholder = "(?)"
	}

	whereClause := fmt.Sprintf("%s %s %s", column, targetOperator, placeholder)
	q.WhereList = append(q.WhereList, whereClause)
	q.BindingList = append(q.BindingList, targetValue)

	return q
}

func (q *QueryBuilder) OrderBy(column string, direction ...string) *QueryBuilder {
	dir := "ASC"
	if len(direction) > 0 {
		dir = strings.ToUpper(direction[0])
	}

	if dir != "ASC" && dir != "DESC" {
		log.Printf("Invalid order direction: %s", dir)
		return q
	}

	if !strings.Contains(column, ".") {
		column = fmt.Sprintf("`%s`", column)
	}

	orderClause := fmt.Sprintf("%s %s", column, dir)
	q.OrderList = append(q.OrderList, orderClause)
	return q
}

func (q *QueryBuilder) Limit(num int) *QueryBuilder {
	q.QueryLimit = &num
	return q
}

func (q *QueryBuilder) Offset(num int) *QueryBuilder {
	q.QueryOffset = &num
	return q
}

func (q *QueryBuilder) Increase(target string, number ...int) *QueryBuilder {
	num := 1
	if len(number) > 0 {
		num = number[0]
	}

	setClause := fmt.Sprintf("%s = %s + %d", target, target, num)
	q.SetList = append(q.SetList, setClause)
	return q
}

// * private method
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
