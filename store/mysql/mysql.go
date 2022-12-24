package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hhzhhzhhz/gopkg/config"
	"github.com/jmoiron/sqlx"
	"net/url"
)

// OptionMysql eg: user:password@tcp(ip:port)/db_name?charset=utf8&parseTime=True&loc=%s&readTimeout=10s&timeout=30s
type Cfg struct {
	Dsn string `toml:"dsn"`
}

func RawConfig(name string) *Cfg {
	var cfg *Cfg
	key := "store." + name
	if err := config.GetConfig().UnmarshalKey(key, &cfg); err != nil {
		panic(fmt.Sprintf("unmarshal key=%s failed cause=%s", key, err.Error()))
	}
	return cfg
}

func (m *Cfg) Validate() error {
	if m.Dsn == "" {
		return fmt.Errorf("mysql dsn is empty")
	}
	return nil
}

func (m *Cfg) Build() (*sqlx.DB, error) {
	dsn := fmt.Sprintf(m.Dsn, url.QueryEscape("Asia/Shanghai"))
	return sqlx.Open("mysql", dsn)
}
