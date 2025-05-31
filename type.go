package golangMysqlPool

import (
	"database/sql"
	"log"
)

type QueryBuilder struct {
	db          *sql.DB
	Database    *string
	TableName   *string
	SelectList  []string
	JoinList    []string
	WhereList   []string
	BindingList []interface{}
	OrderList   []string
	SetList     []string
	QueryLimit  *int
	QueryOffset *int
	WithTotal   bool
	Logger      *Logger
}

type Pool struct {
	db     *sql.DB
	Logger *Logger
}

type PoolList struct {
	Read   *Pool
	Write  *Pool
	Logger *Logger
}

type Config struct {
	Host       string `json:"host,omitempty"`
	Port       int    `json:"port,omitempty"`
	User       string `json:"user,omitempty"`
	Password   string `json:"password,omitempty"`
	Charset    string `json:"charset,omitempty"`
	Connection int    `json:"connection,omitempty"`
}

type ConfigList struct {
	Read    *Config
	Write   *Config
	LogPath string `json:"log_path,omitempty"`
}

type Logger struct {
	InitLogger   *log.Logger
	ActionLogger *log.Logger
	Path         string
}
