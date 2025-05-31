package golangMysqlPool

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func New(c *ConfigList) (*PoolList, error) {
	if c == nil {
		return nil, fmt.Errorf("Config is required")
	}

	if c.LogPath == "" {
		c.LogPath = "./logs/golangMysqlPool"
	}

	logger, err := newLogger(c.LogPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to init logger: %v", err)
	}

	var pool = &PoolList{
		Read:   nil,
		Write:  nil,
		Logger: logger,
	}

	readConfig := c.Read

	if readConfig.Host == "" {
		readConfig.Host = "localhost"
	}

	if readConfig.Port == 0 {
		readConfig.Port = 3306
	}

	if readConfig.User == "" {
		readConfig.User = "root"
	}

	if readConfig.Password == "" {
		readConfig.Password = ""
	}

	if readConfig.Charset == "" {
		readConfig.Charset = "utf8mb4"
	}

	if readConfig.Connection == 0 {
		readConfig.Connection = 4
	}

	read, err := sql.Open("mysql",
		fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/?charset=%s&parseTime=true",
			readConfig.User,
			readConfig.Password,
			readConfig.Host,
			readConfig.Port,
			readConfig.Charset,
		),
	)
	if err != nil {
		logger.Init(true, "Failed to create read pool", err.Error())
		return nil, fmt.Errorf("Failed to create read pool: %w", err)
	}

	read.SetMaxOpenConns(readConfig.Connection)
	read.SetMaxIdleConns(readConfig.Connection / 2)
	read.SetConnMaxLifetime(time.Hour)

	if err := read.Ping(); err != nil {
		logger.Init(true, "Failed to connect read pool", err.Error())
		return nil, fmt.Errorf("Failed to connect read pool: %w", err)
	}

	pool.Read = &Pool{db: read}

	writeConfig := c.Write
	if writeConfig == nil {
		writeConfig = readConfig
	}

	writeDB, err := sql.Open(
		"mysql",
		fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/?charset=%s&parseTime=true",
			writeConfig.User,
			writeConfig.Password,
			writeConfig.Host,
			writeConfig.Port,
			writeConfig.Charset,
		),
	)
	if err != nil {
		logger.Init(true, "Failed to create write pool", err.Error())
		return nil, fmt.Errorf("Failed to create write pool: %w", err)
	}

	writeDB.SetMaxOpenConns(writeConfig.Connection)
	writeDB.SetMaxIdleConns(writeConfig.Connection / 2)
	writeDB.SetConnMaxLifetime(time.Hour)

	if err := writeDB.Ping(); err != nil {
		logger.Init(true, "Failed to connect write pool", err.Error())
		return nil, fmt.Errorf("Failed to connect write pool: %w", err)
	}

	pool.Write = &Pool{db: writeDB}

	pool.listenShutdownSignal()
	pool.Write.Logger = logger
	pool.Read.Logger = logger
	return pool, nil
}

func (p *PoolList) Close() error {
	var readErr, writeErr error

	if p.Read != nil && p.Read.db != nil {
		readErr = p.Read.db.Close()
		p.Read = nil
	}

	if p.Write != nil && p.Write.db != nil {
		writeErr = p.Write.db.Close()
		p.Write = nil
	}

	if readErr != nil {
		p.Write.Logger.Action(true, "Failed to close read pool", readErr.Error())
		return fmt.Errorf("Failed to close read pool: %w", readErr)
	}
	if writeErr != nil {
		p.Write.Logger.Action(true, "Failed to close write pool", writeErr.Error())
		return fmt.Errorf("Failed to close write pool: %w", writeErr)
	}

	return nil
}

// * private method
func (p *PoolList) listenShutdownSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		_ = p.Close()
		os.Exit(0)
	}()
}
