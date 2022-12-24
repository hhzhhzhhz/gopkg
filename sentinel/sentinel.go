package sentinel

import (
	"github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/core/hotspot"
	"github.com/alibaba/sentinel-golang/core/isolation"
	"github.com/alibaba/sentinel-golang/core/stat"
	"github.com/alibaba/sentinel-golang/core/system"
	"github.com/alibaba/sentinel-golang/logging"
	"sync"
	"time"
)

const (
	resetTick = 1 * time.Hour
)

var (
	once sync.Once
)

func InitSentinel(cfg *config.Entity) error {
	if cfg != nil {
		return api.InitWithConfig(cfg)
	}
	conf := config.NewDefaultConfig()
	conf.Sentinel.Log.Logger = logging.NewConsoleLogger()
	return api.InitWithConfig(conf)
}

// Reset 释放资源
func Reset() {
	once.Do(func() {
		t := time.NewTicker(resetTick)
		for {
			select {
			case <-t.C:
				stat.ResetResourceNodeMap()
			}
		}
	})
}

// Entry
// Hotspot: api.Entry(resource, api.WithArgs(args...))
// Isolation: api.Entry(resource, api.WithBatchCount(uint32))
// Breaker: api.Entry(resource)
// Limit: api.Entry(resource, api.WithTrafficType(base.TrafficType))
func Entry(resource string, opts ...api.EntryOption) (*base.SentinelEntry, *base.BlockError) {
	return api.Entry(resource, opts...)
}

func LoadHostRule(rules []*hotspot.Rule) (bool, error) {
	return hotspot.LoadRules(rules)
}

func LoadFlowRule(rules []*flow.Rule) (bool, error) {
	return flow.LoadRules(rules)
}

func LoadCircuitBreakerRule(rules []*circuitbreaker.Rule) (bool, error) {
	return circuitbreaker.LoadRules(rules)
}

func LoadIsolationRule(rules []*isolation.Rule) (bool, error) {
	return isolation.LoadRules(rules)
}

func ClearHotspotRule() error {
	return hotspot.ClearRules()
}

func ClearFlowRule() error {
	return flow.ClearRules()
}

func ClearIsolationRule() error {
	return isolation.ClearRules()
}

func ClearSystemRule() error {
	return system.ClearRules()
}

func ClearCircuitBreakerRule() error {
	return circuitbreaker.ClearRules()
}
