package log

import (
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
	"os"
)

// Level log level
type Level int32

// Logger level
const (
	LvlDebug Level = iota
	LvlInfo
	LvlWarn
	LvlError
	LvlFatal
)

type Log interface {
	Info(format string, v ...interface{})
	Debug(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
	Fatal(format string, v ...interface{})
	Level() Level
	Close() error
}

func Logger() Log {
	return stdLogger
}

// Logger stores a logger
type logger struct {
	log.Logger
	level Level
}

var stdLogger = New(os.Stdout)

// SetOutput sets the writer of standard logger
func SetOutput(w io.Writer) {
	stdLogger.SetOutput(w)
}

// SetLevel set log level
func SetLevel(l Level) {
	stdLogger.SetLevel(l)
}

// GetLevel returns current log level
func GetLevel() Level {
	return stdLogger.Level()
}

func Fatal(format string, v ...interface{}) {
	stdLogger.Output(2, fmt.Sprintf("[FATAL] "+format+"\n", v...))
	os.Exit(1)
}

func Error(format string, v ...interface{}) {
	if stdLogger.Level() <= LvlError {
		stdLogger.Output(2, fmt.Sprintf("[ERROR] "+format+"\n", v...))
	}
}

func Warn(format string, v ...interface{}) {
	if stdLogger.Level() <= LvlWarn {
		stdLogger.Output(2, fmt.Sprintf("[WARN] "+format+"\n", v...))
	}
}

func Info(format string, v ...interface{}) {
	if stdLogger.Level() <= LvlInfo {
		stdLogger.Output(2, fmt.Sprintf("[INFO] "+format+"\n", v...))
	}
}

func Debug(format string, v ...interface{}) {
	if stdLogger.Level() <= LvlDebug {
		stdLogger.Output(3, fmt.Sprintf("[DEBUG] "+format+"\n", v...))
	}
}

// New creates a instance of Logger
func New(w io.Writer) *logger {
	l := &logger{level: LvlInfo}
	l.SetOutput(w)
	l.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmicroseconds)

	return l
}

func (l *logger) Debug(format string, v ...interface{}) {
	if l.level <= LvlDebug {
		l.Output(2, fmt.Sprintf("[DEBUG] "+format+"\n", v...))
	}
}

func (l *logger) Info(format string, v ...interface{}) {
	if l.level <= LvlInfo {
		l.Output(2, fmt.Sprintf("[INFO] "+format+"\n", v...))
	}
}

func (l *logger) Warn(format string, v ...interface{}) {
	if l.level <= LvlWarn {
		l.Output(2, fmt.Sprintf("[WARN] "+format+"\n", v...))
	}
}

func (l *logger) Error(format string, v ...interface{}) {
	if l.level <= LvlError {
		l.Output(2, fmt.Sprintf("[ERROR] "+format+"\n", v...))
	}
}

func (l *logger) Fatal(format string, v ...interface{}) {
	l.Output(2, fmt.Sprintf("[FATAL] "+format+"\n", v...))
	os.Exit(1)
}

func (l *logger) Close() error {
	return nil
}

// Level returns current logger level
func (l *logger) Level() Level {
	return l.level
}

// SetLevel sets the logger level
func (l *logger) SetLevel(level Level) {
	l.level = level
}

// InitLumberjack info、debug、warn、error
func InitLumberjack(level, fileName string) {
	if fileName == "" {
		fileName = "lumberjack"
	}
	writer := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    100,
		MaxBackups: 1,
		MaxAge:     10,
	}
	stdLogger.SetOutput(writer)
	stdLogger.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmicroseconds)
	switch level {
	case "debug", "DEBUG":
		stdLogger.SetLevel(LvlDebug)
	case "warn", "WARN":
		stdLogger.SetLevel(LvlWarn)
	case "error", "ERROR":
		stdLogger.SetLevel(LvlError)
	default:
		stdLogger.SetLevel(LvlInfo)
	}
}

func InitLog(level string, writer io.Writer) {
	stdLogger.SetOutput(writer)
	stdLogger.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmicroseconds)
	switch level {
	case "debug", "DEBUG":
		stdLogger.SetLevel(LvlDebug)
	case "warn", "WARN":
		stdLogger.SetLevel(LvlWarn)
	case "error", "ERROR":
		stdLogger.SetLevel(LvlError)
	default:
		stdLogger.SetLevel(LvlInfo)
	}
}
