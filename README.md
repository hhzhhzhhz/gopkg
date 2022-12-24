# Introduction
 * gopkg 服务轻巧快速开发工具包, 配置加载方式借鉴 jupiter


# Features
 * 配置
 * 选举、分布式锁
 * go 携程池
 * 消息队列nsq、rabbit
 * 日志
 * msyql、gorm、etcd
 * sentinel 限流
 * cron 任务调度
 * singleflight
 * 常用util 工具等等

# More Example
* <a href="./example/" alt="链接">Example(example/**)</a>

# Load Custom Components
```go
type Cfg struct {
	addr string `toml:"addr"`
}

func RawConfig(name string) *Cfg {
	var cfg *Cfg
	// 归属于哪个类别 客户端client.*** 存储store.***
	key := "client." + name
	if err := config.GetConfig().UnmarshalKey(key, &cfg); err != nil {
		panic(fmt.Sprintf("unmarshal key=%s failed cause=%s", key, err.Error()))
	}
	return cfg
}

func (c *Cfg) Validate() error {
	return nil
}

// 通过配置构建实例
func (c *Cfg) Build() (interface, error) {
	return Newxxx(c)
}

// 创建实例
func Newxxx() (interface, error) {
	return nil, nil
}

// 使用
func init() {
    c, err := RawConfig("xxx").Build()
    if err != nil {
        t.Errorf(err.Error())
    }
    c.xxx()
}
```
