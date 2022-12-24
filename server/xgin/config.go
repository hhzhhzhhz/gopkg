package xgin

import (
	"fmt"
	"github.com/hhzhhzhhz/gopkg/config"
)

type Cfg struct {
	Port       int    `toml:"port"`
	StaticPath string `toml:"static_path"`
	StaticFS   string `toml:"static_fs"`
}

func RawConfig(name string) *Cfg {
	var cfg *Cfg
	key := "server." + name
	if err := config.GetConfig().UnmarshalKey(key, &cfg); err != nil {
		panic(fmt.Sprintf("unmarshal key=%s failed cause=%s", key, err.Error()))
	}
	return cfg
}

func (m *Cfg) Validate() error {
	if m.Port == 0 {
		return fmt.Errorf("gin port is empty")
	}
	return nil
}

func (m *Cfg) Build() (*Server, error) {
	return NewServer(m), nil
}
