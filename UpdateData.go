package golangMysqlPool

import (
	"database/sql"
	"fmt"
	"strings"
)

func (q *QueryBuilder) Update(data ...map[string]interface{}) (sql.Result, error) {
	if q.TableName == nil {
		q.Logger.Action(true, "Table is required")
		return nil, fmt.Errorf("Table is required")
	}

	values := []interface{}{}

	if len(data) > 0 {
		for column, value := range data[0] {
			columnName := column
			if !strings.Contains(column, ".") {
				columnName = fmt.Sprintf("`%s`", column)
			}

			if str, ok := value.(string); ok && contains(supportFunction, strings.ToUpper(str)) {
				q.SetList = append(q.SetList, fmt.Sprintf("%s = %s", columnName, str))
			} else {
				q.SetList = append(q.SetList, fmt.Sprintf("%s = ?", columnName))
				values = append(values, value)
			}
		}
	}

	query := fmt.Sprintf("UPDATE `%s` SET %s", *q.TableName, strings.Join(q.SetList, ", "))

	if len(q.WhereList) > 0 {
		query += " WHERE " + strings.Join(q.WhereList, " AND ")
	}

	allValues := append(values, q.BindingList...)
	return q.exec(query, allValues...)
}
