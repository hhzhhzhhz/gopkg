package kube

import (
	"fmt"
	"github.com/hhzhhzhhz/gopkg/config"
	kubeclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

type Cfg struct {
	// 配置路径
	Path string `toml:"path"`
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

func (c *Cfg) Build() (*kubeclient.Clientset, error) {
	return NewKubeCli(c)
}

func NewKubeCli(opt *Cfg) (*kubeclient.Clientset, error) {
	var kluge *rest.Config
	var err error
	if opt.Path == "" {
		kluge, err = rest.InClusterConfig()
		if err != nil {
			opt.Path = filepath.Join(os.Getenv("HOME"), "/etc", "config")
			kluge, err = clientcmd.BuildConfigFromFlags("", opt.Path)
			if err != nil {
				return nil, fmt.Errorf("BuildConfigFromFlags kluge failed cause=%s", err.Error())
			}
		}
	} else {
		kluge, err = clientcmd.BuildConfigFromFlags("", opt.Path)
		if err != nil {
			return nil, fmt.Errorf("BuildConfigFromFlags kluge failed cause=%s", err.Error())
		}
	}
	return kubeclient.NewForConfig(kluge)
}
