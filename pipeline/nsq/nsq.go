package nsq

import (
	"encoding/json"
	"fmt"
	"github.com/hhzhhzhhz/gopkg/config"
	"github.com/nsqio/go-nsq"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	nsqInterval = 10
	nsqNodesApi = "http://%s/nodes"
)

type Cfg struct {
	LoopAddr string `toml:"loop_addr"`
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
	if m.LoopAddr == "" {
		return fmt.Errorf("nsq loop_addr is empty")
	}
	return nil
}

func (m *Cfg) Build() (MQ, error) {
	return NewNsq(m)
}

type MQ interface {
	Publish(topic string, message interface{}) error
	Subscription(topic, channel string, handler nsq.Handler) error
	UnSubscription(topic string)
	Close() error
}

type ns struct {
	*Cfg
	mux        sync.Mutex
	producer   *nsq.Producer
	consumers  []*Consumers
	nsqdAdders []string
}

type Consumers struct {
	topic   string
	channel string
	handler nsq.Handler
	c       *nsq.Consumer
}

func NewNsq(opt *Cfg) (MQ, error) {
	mq := &ns{Cfg: opt}
	if err := mq.nsqdParse(mq.LoopAddr); err != nil {
		return nil, err
	}
	if err := mq.initProducer(); err != nil {
		return nil, err
	}
	return mq, nil
}

func (n *ns) initProducer() error {
	p, err := nsq.NewProducer(n.nsqdAdders[0], nsq.NewConfig())
	if err != nil {
		return err
	}
	p.SetLoggerLevel(nsq.LogLevelError)
	n.producer = p
	return nil
}

func (n *ns) initConsumers(cs []*Consumers) error {
	if len(cs) == 0 {
		return nil
	}
	for _, c := range cs {
		if err := n.subscription(c.topic, c.channel, c.handler); err != nil {
			return err
		}
	}
	return nil
}

// todo build connect pool
func (n *ns) nsqdParse(lookup string) error {
	resp, err := http.Get(fmt.Sprintf(nsqNodesApi, lookup))
	if err != nil {
		return err
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()
	endpoint, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	lupresp := &LookupResp{}
	if err := json.Unmarshal(endpoint, lupresp); err != nil {
		return err
	}
	var nsqdAdders []string
	for _, producer := range lupresp.Producers {
		broadcastAddress := producer.BroadcastAddress
		port := producer.TCPPort
		joined := net.JoinHostPort(broadcastAddress, strconv.Itoa(port))
		nsqdAdders = append(nsqdAdders, joined)
	}
	if len(nsqdAdders) == 0 {
		return fmt.Errorf("nsqd addrs is empty")
	}
	n.nsqdAdders = nsqdAdders
	return nil
}

func (n *ns) Publish(topic string, message interface{}) error {
	n.mux.Lock()
	defer n.mux.Unlock()
	return n.publish(topic, message)
}

func (n *ns) publish(topic string, message interface{}) error {
	var b []byte
	b, ok := message.([]byte)
	if !ok {
		var err error
		b, err = json.Marshal(message)
		if err != nil {
			return fmt.Errorf("topic: %s json.marshal failed cause: %s", topic, err.Error())
		}
	}
	if err := n.producer.Publish(topic, b); err != nil {
		return fmt.Errorf("topic: %s cause: %s", topic, err.Error())
	}
	return nil
}

func (n *ns) Subscription(topic, channel string, handler nsq.Handler) error {
	n.mux.Lock()
	defer n.mux.Unlock()
	return n.subscription(topic, channel, handler)
}

func (n *ns) subscription(topic, channel string, handler nsq.Handler) error {
	if channel == "" {
		channel = topic
	}
	cfg := nsq.NewConfig()
	cfg.MaxInFlight = 10
	cfg.LookupdPollInterval = nsqInterval * time.Second
	c, err := nsq.NewConsumer(topic, channel, cfg)
	if err != nil {
		return err
	}
	c.SetLoggerLevel(nsq.LogLevelError)
	c.AddHandler(handler)
	err = c.ConnectToNSQLookupd(n.LoopAddr)
	if err != nil {
		return err
	}

	n.consumers = append(n.consumers, &Consumers{topic: topic, c: c})
	return nil
}

func (n *ns) UnSubscription(topic string) {
	n.mux.Lock()
	defer n.mux.Unlock()
	for _, v := range n.consumers {
		if topic == v.topic {
			v.c.Stop()
		}
	}
}

func (n *ns) Close() error {
	n.mux.Lock()
	defer n.mux.Unlock()
	var errs error
	if n.producer != nil {
		n.producer.Stop()
	}
	for _, c := range n.consumers {
		c.c.Stop()
	}
	return errs
}

type LookupResp struct {
	Producers []*peerInfo `json:"producers"`
}

type peerInfo struct {
	RemoteAddress    string `json:"remote_address"`
	Hostname         string `json:"hostname"`
	BroadcastAddress string `json:"broadcast_address"`
	TCPPort          int    `json:"tcp_port"`
	HTTPPort         int    `json:"http_port"`
	Version          string `json:"version"`
}
