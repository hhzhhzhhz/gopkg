package sentinel

import (
	"fmt"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/core/hotspot"
	"github.com/alibaba/sentinel-golang/core/isolation"
	"github.com/hhzhhzhhz/gopkg/config"
)

type Cfg struct {
	Hr []*hotspot.Rule
	Fr []*flow.Rule
	Cr []*circuitbreaker.Rule
	Ir []*isolation.Rule
}

func (e *Cfg) Validate() error {
	return nil
}

func (e *Cfg) BuildHr() (bool, error) {
	return LoadHostRule(e.Hr)
}

func (e *Cfg) BuildFr() (bool, error) {
	return LoadFlowRule(e.Fr)
}

func (e *Cfg) BuildIr() (bool, error) {
	return LoadIsolationRule(e.Ir)
}

func (e *Cfg) BuildCr() (bool, error) {
	return LoadCircuitBreakerRule(e.Cr)
}

func RawConfig(name string) *Cfg {
	cfg := &Cfg{Hr: []*hotspot.Rule{}, Fr: []*flow.Rule{}, Cr: []*circuitbreaker.Rule{}, Ir: []*isolation.Rule{}}
	key := "sentinel." + name
	if err := config.GetConfig().UnmarshalKey(key, &cfg); err != nil {
		panic(fmt.Sprintf("unmarshal key=%s failed cause=%s", key, err.Error()))
	}
	return cfg
}
