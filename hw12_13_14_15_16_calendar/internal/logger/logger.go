package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

type logLevel int

const (
	DEBUG logLevel = iota
	INFO
	WARNING
	ERROR
)

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warning(msg string)
	Error(msg string)
}

type Logg struct {
	out   io.Writer
	level logLevel
	mu    sync.Mutex
}

var levelsToInt = map[string]logLevel{
	"DEBUG":   DEBUG,
	"INFO":    INFO,
	"WARNING": WARNING,
	"ERROR":   ERROR,
}

var intToLevels = map[logLevel]string{
	DEBUG:   "DEBUG",
	INFO:    "INFO",
	WARNING: "WARNING",
	ERROR:   "ERROR",
}

func New(level string) *Logg {
	logger := &Logg{}
	logger.SetLevel(level)
	logger.SetOutput(os.Stdout)

	return logger
}

func (l *Logg) SetLevel(level string) {
	curLevel, ok := levelsToInt[level]
	if !ok {
		log.Fatal(fmt.Errorf("invalid log level: %s", level))
	}

	l.level = curLevel
}

func (l *Logg) SetOutput(out io.Writer) {
	l.out = out
}

func (l *Logg) Debug(msg string) {
	l.log(DEBUG, msg)
}

func (l *Logg) Info(msg string) {
	l.log(INFO, msg)
}

func (l *Logg) Warning(msg string) {
	l.log(WARNING, msg)
}

func (l *Logg) Error(msg string) {
	l.log(ERROR, msg)
}

func (l *Logg) log(level logLevel, msg string) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	msg = fmt.Sprintf("%s %s %s\n", time.Now().Format(time.DateTime), intToLevels[level], msg)

	_, err := l.out.Write([]byte(msg))
	if err != nil {
		log.Printf("error writing to log out: %v\n", err)
	}
}
