package etcd

import (
	"fmt"
	"github.com/hhzhhzhhz/gopkg/config"
	"go.etcd.io/etcd/client/v3"
	"time"
)

const (
	etcdDefaultTimeout = 5 * time.Second
)

type Cfg struct {
	*clientv3.Config
}

func (c *Cfg) Validate() error {
	if len(c.Endpoints) == 0 {
		return fmt.Errorf("ectd endpoints is empty")
	}
	return nil
}

func RawConfig(name string) *Cfg {
	cfg := &Cfg{Config: &clientv3.Config{}}
	key := "store." + name
	if err := config.GetConfig().UnmarshalKey(key, &cfg.Config); err != nil {
		panic(fmt.Sprintf("unmarshal key=%s failed cause=%s", key, err.Error()))
	}
	return cfg
}

func (c *Cfg) Build() (*clientv3.Client, error) {
	return clientv3.New(*c.Config)
}
