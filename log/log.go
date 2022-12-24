package log

import (
	"github.com/hhzhhzhhz/logs-go"
	"os"
	"sync"
)

var (
	onceJ   sync.Once
	stdLogJ logs_go.LogJ
	stdLogF logs_go.Logf
)

func SetLogger(log logs_go.Logf) {
	stdLogF = log
}

func Logger() logs_go.Logf {
	if stdLogF != nil {
		return stdLogF
	}
	return logs_go.DefalutLogf()
}

func LoggerJ() logs_go.LogJ {
	onceJ.Do(func() {
		hostname, _ := os.Hostname()
		cfg := logs_go.NewLogJconfig()
		cfg.Stdout = true
		cfg.InitialFields = map[string]interface{}{"hostname": hostname}
		stdLogJ, _ = cfg.BuildLogJ()
	})
	return stdLogJ
}
