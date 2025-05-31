package golangMysqlPool

import (
	"database/sql"
	"fmt"
	"time"
)

func (db *Pool) Query(query string, params ...interface{}) (*sql.Rows, error) {
	if db.db == nil {
		db.Logger.Action(true, "Database connection is not available")
		return nil, fmt.Errorf("database connection is not available")
	}

	startTime := time.Now()
	rows, err := db.db.Query(query, params...)
	duration := time.Since(startTime)

	if duration > 20*time.Millisecond {
		db.Logger.Action(false,
			fmt.Sprintf("Slow Query: %v", duration),
			query,
		)
	}

	return rows, err
}

func (db *Pool) Exec(query string, params ...interface{}) (sql.Result, error) {
	if db.db == nil {
		db.Logger.Action(true, "Database connection is not available")
		return nil, fmt.Errorf("database connection is not available")
	}

	startTime := time.Now()
	result, err := db.db.Exec(query, params...)
	duration := time.Since(startTime)

	if duration > 20*time.Millisecond {
		db.Logger.Action(false,
			fmt.Sprintf("Slow Query: %v", duration),
			query,
		)
	}

	return result, err
}

// * private method
func (q *QueryBuilder) query(query string, params ...interface{}) (*sql.Rows, error) {
	if q.db == nil {
		q.Logger.Action(true, "Database connection is not available")
		return nil, fmt.Errorf("database connection is not available")
	}

	startTime := time.Now()
	rows, err := q.db.Query(query, params...)
	duration := time.Since(startTime)

	if duration > 20*time.Millisecond {
		q.Logger.Action(false,
			fmt.Sprintf("Slow Query: %v", duration),
			query,
		)
	}

	return rows, err
}

// * private method
func (q *QueryBuilder) exec(query string, params ...interface{}) (sql.Result, error) {
	if q.db == nil {
		q.Logger.Action(true, "Database connection is not available")
		return nil, fmt.Errorf("database connection is not available")
	}

	startTime := time.Now()
	result, err := q.db.Exec(query, params...)
	duration := time.Since(startTime)

	if duration > 20*time.Millisecond {
		q.Logger.Action(false,
			fmt.Sprintf("Slow Query: %v", duration),
			query,
		)
	}

	return result, err
}
