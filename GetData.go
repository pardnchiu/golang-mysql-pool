package golangMysqlPool

import (
	"database/sql"
	"fmt"
	"strings"
)

func (q *QueryBuilder) Get() (*sql.Rows, error) {
	if q.TableName == nil {
		q.Logger.Action(true, "Table is required")
		return nil, fmt.Errorf("Table is required")
	}

	fieldNames := make([]string, len(q.SelectList))
	for i, field := range q.SelectList {
		switch {
		case field == "*":
			fieldNames[i] = "*"
		case strings.ContainsAny(field, ".()"):
			fieldNames[i] = field
		default:
			fieldNames[i] = fmt.Sprintf("`%s`", field)
		}
	}

	query := fmt.Sprintf("SELECT %s FROM `%s`", strings.Join(fieldNames, ", "), *q.TableName)

	if len(q.JoinList) > 0 {
		query += " " + strings.Join(q.JoinList, " ")
	}

	if len(q.WhereList) > 0 {
		query += " WHERE " + strings.Join(q.WhereList, " AND ")
	}

	if q.WithTotal {
		query = fmt.Sprintf("SELECT COUNT(*) OVER() AS total, data.* FROM (%s) AS data", query)
	}

	if len(q.OrderList) > 0 {
		query += " ORDER BY " + strings.Join(q.OrderList, ", ")
	}

	if q.QueryLimit != nil {
		query += fmt.Sprintf(" LIMIT %d", *q.QueryLimit)
	}

	if q.QueryOffset != nil {
		query += fmt.Sprintf(" OFFSET %d", *q.QueryOffset)
	}

	return q.query(query, q.BindingList...)
}
