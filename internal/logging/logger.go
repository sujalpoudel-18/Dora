package logging

import (
	"fmt"
	"io"
	"os"
	"sync"
)

type Level int

const (
	LevelSilent Level = iota
	LevelInfo
	LevelVerbose
	LevelDebug
)

type Logger struct {
	mu     sync.Mutex
	level  Level
	writer io.Writer
}

var defaultLogger = NewLogger(LevelInfo, os.Stdout)

func NewLogger(level Level, writer io.Writer) *Logger {
	return &Logger{
		level:  level,
		writer: writer,
	}
}

func SetLevel(level Level) {
	defaultLogger.mu.Lock()
	defer defaultLogger.mu.Unlock()
	defaultLogger.level = level
}

func Info(format string, v ...interface{}) {
	defaultLogger.log(LevelInfo, "[*] ", format, v...)
}

func Verbose(format string, v ...interface{}) {
	defaultLogger.log(LevelVerbose, "[+] ", format, v...)
}

func Debug(format string, v ...interface{}) {
	defaultLogger.log(LevelDebug, "[DEBUG] ", format, v...)
}

func Error(format string, v ...interface{}) {
	defaultLogger.log(LevelSilent, "[-] ", format, v...)
}

func (l *Logger) log(level Level, prefix, format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if level > l.level && level != LevelSilent {
		return
	}

	msg := fmt.Sprintf(format, v...)
	fmt.Fprintf(l.writer, "%s%s\n", prefix, msg)
}
