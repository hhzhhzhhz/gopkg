package nsq_admin

import (
	"context"
	"fmt"
	"github.com/hhzhhzhhz/gopkg/config"
	"github.com/hhzhhzhhz/gopkg/utils"
	"net/http"
)

type Cfg struct {
	AdminAddr string `toml:"admin_addr"`
}

func RawConfig(name string) *Cfg {
	var cfg *Cfg
	key := "client." + name
	if err := config.GetConfig().UnmarshalKey(key, &cfg); err != nil {
		panic(fmt.Sprintf("unmarshal key=%s failed cause=%s", key, err.Error()))
	}
	return cfg
}

func (c *Cfg) Validate() error {
	return nil
}

func (c *Cfg) Build() (*NsqSeq, error) {
	return &NsqSeq{Cfg: c}, nil
}

type NsqSeq struct {
	*Cfg
}

func (n *NsqSeq) DeleteChannel(ctx context.Context, channel string) error {
	url := fmt.Sprintf("%s/api/topics/%s", n.AdminAddr, channel)
	request, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("DeleteChannel NewRequest failed cause=%s", err.Error())
	}
	ret := make(map[string]string)
	if err := utils.PostDo(request, &ret); err != nil {
		return fmt.Errorf("DeleteChannel PostDo failed cause=%s", err.Error())
	}
	return nil
}
