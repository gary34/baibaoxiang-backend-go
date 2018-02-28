package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type Logger struct {
	writer io.Writer
	stack  int
}

func NewStdLogger() *Logger {
	return &Logger{
		writer: os.Stdout,
		stack:  2,
	}
}

var gLogger *Logger = NewStdLogger()

const kFileWidth = 16

type LogLevel uint32

const (
	LogLevelError = LogLevel(iota)
	LogLevelWarning
	LogLevelNotice
	LogLevelInfo
	LogLevelData
	LogLevelDebug
	LogLevelFatal
)

var LogLevelTag = map[LogLevel]string{
	LogLevelError:   "Error",
	LogLevelWarning: "Warn",
	LogLevelNotice:  "Notice",
	LogLevelInfo:    "Info",
	LogLevelDebug:   "Debug",
	LogLevelData:    "Data",
	LogLevelFatal:   "Fatal",
}

func CastLogLevel2Tag(level LogLevel) string {
	if item, found := LogLevelTag[level]; found {
		return item
	}

	return "Unknown"
}

func (l *Logger) record(level LogLevel, message string) {
	_, file, line, _ := runtime.Caller(l.stack)
	basename := filepath.Base(file)
	if len(basename) > kFileWidth {
		basename = ".." + basename[len(basename)-kFileWidth+2:]
	}
	out := fmt.Sprintf(
		"%s %14s:%-4d[%-6s] %s\n",
		time.Now().Format("2006-01-02 15:04:05.000"), basename, line, CastLogLevel2Tag(level), message,
	)
	l.writer.Write([]byte(out))
}

func (l *Logger) ILog(format string, args ...interface{}) {
	l.record(LogLevelInfo, fmt.Sprintf(format, args...))
}
func (l *Logger) DLog(format string, args ...interface{}) {
	l.record(LogLevelDebug, fmt.Sprintf(format, args...))
}

func (l *Logger) WLog(format string, args ...interface{}) {
	l.record(LogLevelWarning, fmt.Sprintf(format, args...))
}

func (l *Logger) ELog(format string, args ...interface{}) {
	l.record(LogLevelError, fmt.Sprintf(format, args...))
}
