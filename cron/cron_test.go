package cron

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"testing"
	"time"
)

func TestNewCron(t *testing.T) {
	c := NewCron(cron.WithSeconds(), cron.WithLogger(&logTest{}))
	c.AddDelayJob(10*time.Second, func() {
		t.Log("xxx")
	})
	c.Run()

}

type logTest struct {
}

func (l *logTest) Info(msg string, keysAndValues ...interface{}) {
	fmt.Println("INFO", msg, keysAndValues)
}

func (l *logTest) Error(err error, msg string, keysAndValues ...interface{}) {
	fmt.Println("ERROR", err.Error(), msg, keysAndValues)

}

// Error logs an error condition.
