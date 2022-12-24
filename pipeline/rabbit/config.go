package rabbit

import (
	"fmt"
	"github.com/hhzhhzhhz/gopkg/config"
)

type Cfg struct {
	Url string
}

func RawConfig(name string) *Cfg {
	var cfg *Cfg
	key := "pipeline." + name
	if err := config.GetConfig().UnmarshalKey(key, &cfg); err != nil {
		panic(fmt.Sprintf("unmarshal key=%s failed cause=%s", key, err.Error()))
	}
	return cfg
}

func (m *Cfg) Validate() error {
	if m.Url == "" {
		return fmt.Errorf("rabbit url is empty")
	}
	return nil
}

func (m *Cfg) BuildDelay() (DelayQueue, error) {
	return NewDelayQueue(m.Url), nil
}

func (m *Cfg) BuildCommon() (CommonQueue, error) {
	return NewCommonQueue(m.Url), nil
}
