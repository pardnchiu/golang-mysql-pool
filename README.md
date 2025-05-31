# MySQL Pool (Golang)

> A high-performance MySQL connection pool implementation focusing on chainable syntax, high readability, native methods, and ultra-fast query performance.<br>
> Features read-write separation, query builder, comprehensive logging, and automatic resource management.<br>
> version Node.js can get [here](https://github.com/pardnchiu/nodejs-mysql-pool)<br>
> version PHP can get [here](https://github.com/pardnchiu/php-mysql-pool)

[![version](https://img.shields.io/github/v/tag/pardnchiu/golang-mysql-pool)](https://github.com/pardnchiu/golang-mysql-pool)

## Feature

- ### Dual Connection Pool System
  - Read Pool (optimized for `SELECT`)
  - Write Pool (optimized for `INSERT`/`UPDATE`/`DELETE`)
  - Independent configuration for each pool
  - Automatic connection lifecycle management
- ### Query Builder
  - Chainable methods for building complex queries
  - Support for `SELECT`, `INSERT`, `UPDATE`, `UPSERT`
  - Prepared statement support for SQL injection prevention
- ### Performance Monitoring
  - Slow query detection (configurable threshold: 20ms)
  - Execution time logging
  - Structured error handling and reporting
- ### Advanced Query Features
  - `JOIN` support (`INNER`, `LEFT`, `RIGHT`)
  - `WHERE` conditions with multiple operators
  - `ORDER` BY with `ASC`/`DESC`
  - `LIMIT` and `OFFSET` pagination
  - `COUNT(*)` `OVER()` for total record counting
- ### Smart Resource Management
  - Shutdown signal handling
  - Automatic connection pool cleanup
  - Connection reuse and optimization
  - Configurable connection limits
- ### MySQL Function Support
  - Built-in function recognition (`NOW()`, `UUID()`, `CURRENT_TIMESTAMP`, etc.)
  - Safe function execution without parameterization
  - Comprehensive function whitelist

## How to use

- ### Installation
  ```bash
  go get github.com/pardnchiu/golang-mysql-pool
  ```
- ### Initialize
  ```go
  package main

  import (
    "log"

    "github.com/pardnchiu/golang-mysql-pool"
  )

  func main() {
    config := &golangMysqlPool.ConfigList{
      Read: &golangMysqlPool.Config{
        Host:       "localhost",
        Port:       3306,
        User:       "root",
        Password:   "password",
        Charset:    "utf8mb4",
        Connection: 10,
      },
      Write: &golangMysqlPool.Config{
        Host:       "localhost",
        Port:       3306,
        User:       "root",
        Password:   "password",
        Charset:    "utf8mb4",
        Connection: 5,
      },
      LogPath: "./logs/mysql-pool",
    }

    pool, err := golangMysqlPool.New(config)
    if err != nil {
      log.Fatal("Failed to initialize pool:", err)
    }
    defer pool.Close()
  }
  ```

### SELECT

```go
func selectData(pool *golangMysqlPool.PoolList) {
  rows, err := pool.Read.DB("myapp").
    Table("users").
    Select("id", "name", "email").
    Where("status", "active").
    Where("age", ">", 18).
    OrderBy("created_at", "DESC").
    Limit(10).
    Offset(0).
    Get()
  
  if err != nil {
    log.Printf("Query error: %v", err)
    return
  }
  defer rows.Close()

  for rows.Next() {
    var id int
    var name, email string
    err := rows.Scan(&id, &name, &email)
    if err != nil {
      log.Printf("Scan error: %v", err)
      continue
    }
    fmt.Printf("User: %d, %s, %s\n", id, name, email)
  }
}
```

### INSERT 

```go
func insertData(pool *golangMysqlPool.PoolList) {
  data := map[string]interface{}{
    "name":       "John Doe",
    "email":      "john@example.com",
    "status":     "active",
    "created_at": "NOW()",
  }

  lastID, err := pool.Write.DB("myapp").
    Table("users").
    Insert(data)
  
  if err != nil {
    log.Printf("Insert error: %v", err)
    return
  }

  fmt.Printf("Inserted user with ID: %d\n", lastID)
}
```

### UPDATE 

```go
func updateData(pool *golangMysqlPool.PoolList) {
  updateData := map[string]interface{}{
    "status":     "inactive",
    "updated_at": "NOW()",
  }

  result, err := pool.Write.DB("myapp").
    Table("users").
    Where("email", "john@example.com").
    Update(updateData)
  
  if err != nil {
    log.Printf("Update error: %v", err)
    return
  }

  rowsAffected, _ := result.RowsAffected()
  fmt.Printf("Updated %d rows\n", rowsAffected)
}
```

### UPSERT 

```go
func upsertData(pool *golangMysqlPool.PoolList) {
  data := map[string]interface{}{
    "email":      "jane@example.com",
    "name":       "Jane Smith",
    "status":     "active",
    "created_at": "NOW()",
  }

  updateData := map[string]interface{}{
    "name":       "Jane Smith Updated",
    "updated_at": "NOW()",
  }

  lastID, err := pool.Write.DB("myapp").
    Table("users").
    Upsert(data, updateData)
  
  if err != nil {
    log.Printf("Upsert error: %v", err)
    return
  }

  fmt.Printf("Upserted user with ID: %d\n", lastID)
}
```

### JOIN

```go
func complexQuery(pool *golangMysqlPool.PoolList) {
  rows, err := pool.Read.DB("myapp").
    Table("users").
    Select("users.name", "profiles.bio", "COUNT(*) OVER() as total").
    LeftJoin("profiles", "users.id", "profiles.user_id").
    Where("users.status", "active").
    Where("profiles.is_public", true).
    OrderBy("users.created_at", "DESC").
    Total().
    Limit(20).
    Get()
  
  if err != nil {
    log.Printf("Complex query error: %v", err)
    return
  }
  defer rows.Close()
}
```

### Direct

```go
func directOperations(pool *golangMysqlPool.PoolList) {
  rows, err := pool.Read.Query("SELECT COUNT(*) FROM users WHERE status = ?", "active")
  if err != nil {
    log.Printf("Direct query error: %v", err)
    return
  }
  defer rows.Close()

  result, err := pool.Write.Exec("UPDATE users SET last_login = NOW() WHERE id = ?", 123)
  if err != nil {
    log.Printf("Direct exec error: %v", err)
    return
  }
  
  rowsAffected, _ := result.RowsAffected()
  fmt.Printf("Updated %d rows\n", rowsAffected)
}
```

## Configuration

### ConfigList
- `Read`: Read database configuration
- `Write`: Write database configuration (optional, defaults to Read config)
- `LogPath`: Log file directory path (default: "./logs/golangMysqlPool")

### Config
- `Host`: Database host (default: "localhost")
- `Port`: Database port (default: 3306)
- `User`: Database username (default: "root")
- `Password`: Database password (default: "")
- `Charset`: Character set (default: "utf8mb4")
- `Connection`: Maximum connections (default: 4)

## Supported Query Methods

### Query Builder Methods
1. **DB(name)**: Select database
2. **Table(name)**: Select table
3. **Select(fields...)**: Specify fields to select
4. **Where(column, operator, value)**: Add WHERE condition
5. **InnerJoin/LeftJoin/RightJoin**: Add JOIN clauses
6. **OrderBy(column, direction)**: Add ORDER BY
7. **Limit(num)**: Set LIMIT
8. **Offset(num)**: Set OFFSET
9. **Total()**: Add COUNT(*) OVER() for pagination
10. **Increase(field, amount)**: Increment field value

### Data Operations
1. **Get()**: Execute SELECT query
2. **Insert(data)**: Insert new record
3. **Update(data)**: Update existing records
4. **Upsert(data, updateData)**: Insert or update on duplicate key
1. **Query(sql, params...)**: Execute raw SELECT query
2. **Exec(sql, params...)**: Execute raw INSERT/UPDATE/DELETE