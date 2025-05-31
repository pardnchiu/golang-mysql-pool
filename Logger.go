package golangMysqlPool

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

// * private method
func newLogger(path string) (*Logger, error) {
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %v", err)
	}

	initLog, initErr := openLog(path, "init.log")
	actionLog, createErr := openLog(path, "action.log")

	if initErr != nil || createErr != nil {
		return nil, fmt.Errorf("failed to open log files: init=%v, action=%v",
			initErr, actionLog)
	}

	return &Logger{
		InitLogger:   log.New(io.MultiWriter(initLog, os.Stdout), "", log.LstdFlags),
		ActionLogger: log.New(io.MultiWriter(actionLog, os.Stdout), "", log.LstdFlags),
		Path:         path,
	}, nil
}

func openLog(path string, target string) (*os.File, error) {
	file, err := os.OpenFile(filepath.Join(path, target), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %v", target, err)
	}
	return file, nil
}

func writeToLog(target *log.Logger, isError bool, message ...string) {
	if len(message) == 0 {
		return
	}

	state := ""
	if isError {
		state = "[ERROR] "
	}
	for i, msg := range message {
		if i == 0 {
			target.Printf("%s%s", state, message[i])
		} else if i == len(message)-1 {
			target.Printf("└── %s", msg)
		} else {
			target.Printf("├── %s", msg)
		}
	}
}

func (l *Logger) Init(isError bool, message ...string) {
	writeToLog(l.InitLogger, isError, message...)
}

func (l *Logger) Action(isError bool, message ...string) {
	writeToLog(l.ActionLogger, isError, message...)
}
