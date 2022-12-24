package gorm

import (
	"fmt"
	"github.com/hhzhhzhhz/gopkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Cfg struct {
	*mysql.Config
}

func (c *Cfg) Validate() error {
	if c.DSN == "" {
		return fmt.Errorf("grom dsn is empty")
	}
	return nil
}

func (c *Cfg) Build() (*gorm.DB, error) {
	return gorm.Open(mysql.New(*c.Config), &gorm.Config{})
}

func RawConfig(name string) *Cfg {
	cfg := &Cfg{Config: &mysql.Config{}}
	key := "store." + name
	if err := config.GetConfig().UnmarshalKey(key, &cfg.Config); err != nil {
		panic(fmt.Sprintf("unmarshal key=%s failed cause=%s", key, err.Error()))
	}
	return cfg
}
