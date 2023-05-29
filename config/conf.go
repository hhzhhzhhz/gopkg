package config

import (
	"errors"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/hhzhhzhhz/gopkg/utils"
	"github.com/hhzhhzhhz/gopkg/utils/xcast"
	"github.com/mitchellh/mapstructure"
	"io/ioutil"
	"log"
	"reflect"
	"strings"
	"sync"
	"time"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", utils.Env("CONFIG", "./config.toml"), "config path eg: ./config.toml")
	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Printf("load config from datasource[%s] failed: %v", configPath, err)
		return
	}
	if err := LoadFromDataSource(b, toml.Unmarshal); err != nil {
		log.Printf("load config from datasource[%s] failed: %v", configPath, err)
	}
}

type Config struct {
	mu       sync.RWMutex
	keyDelim string
	override map[string]interface{}
	keyMap   *sync.Map
}

const (
	defaultKeyDelim = "."
)

func New() *Config {
	return &Config{
		override: make(map[string]interface{}),
		keyDelim: defaultKeyDelim,
		keyMap:   &sync.Map{},
	}
}

// Unmarshaller ...
type Unmarshaller = func([]byte, interface{}) error

func (c *Config) reflush(content []byte, unmarshal Unmarshaller) error {
	configuration := make(map[string]interface{})
	if err := unmarshal(content, &configuration); err != nil {
		return err
	}
	for k, v := range configuration {
		c.Set(k, v)
	}
	return nil
}

func (c *Config) apply(conf map[string]interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	utils.MergeStringMap(c.override, conf)
	for k, v := range c.traverse(c.keyDelim) {
		orig, ok := c.keyMap.Load(k)
		if ok && !reflect.DeepEqual(orig, v) {
		}
		c.keyMap.Store(k, v)
	}

	return nil
}

// Set ...
func (c *Config) Set(key string, val interface{}) error {
	paths := strings.Split(key, c.keyDelim)
	lastKey := paths[len(paths)-1]
	m := deepSearch(c.override, paths[:len(paths)-1])
	m[lastKey] = val
	return c.apply(m)
}

func deepSearch(m map[string]interface{}, path []string) map[string]interface{} {
	for _, k := range path {
		m2, ok := m[k]
		if !ok {
			m3 := make(map[string]interface{})
			m[k] = m3
			m = m3
			continue
		}
		m3, ok := m2.(map[string]interface{})
		if !ok {
			m3 = make(map[string]interface{})
			m[k] = m3
		}
		m = m3
	}
	return m
}

// Get ...
func (c *Config) Get(key string) interface{} {
	return c.find(key)
}


// GetString returns the value associated with the key as a string.
func (c *Config) GetString(key string) string {
	return xcast.ToString(c.Get(key))
}

// GetBool returns the value associated with the key as a boolean.
func (c *Config) GetBool(key string) bool {
	return xcast.ToBool(c.Get(key))
}

func (c *Config) GetInt(key string) int {
	return xcast.ToInt(c.Get(key))
}

// GetInt64 returns the value associated with the key as an integer.
func (c *Config) GetInt64(key string) int64 {
	return xcast.ToInt64(c.Get(key))
}

func (c *Config) GetFloat64(key string) float64 {
	return xcast.ToFloat64(c.Get(key))
}

// GetTime returns the value associated with the key as time.
func (c *Config) GetTime(key string) time.Time {
	return xcast.ToTime(c.Get(key))
}

// GetDuration returns the value associated with the key as a duration.
func (c *Config) GetDuration(key string) time.Duration {
	return xcast.ToDuration(c.Get(key))
}

// GetStringSlice returns the value associated with the key as a slice of strings.
func (c *Config) GetStringSlice(key string) []string {
	return xcast.ToStringSlice(c.Get(key))
}

// GetInt64Slice returns the value associated with the key as a slice of int64s.
func (c *Config) GetInt64Slice(key string) []int64 {
	return xcast.ToInt64Slice(c.Get(key))
}

// GetStringMapString returns the value associated with the key as a map of strings.
func (c *Config) GetStringMapString(key string) map[string]string {
	return xcast.ToStringMapString(c.Get(key))
}

func (c *Config) GetStringMapStringSlice(key string) map[string][]string {
	return xcast.ToStringMapStringSlice(c.Get(key))
}

func (c *Config) find(key string) interface{} {
	dd, ok := c.keyMap.Load(key)
	if ok {
		return dd
	}
	paths := strings.Split(key, c.keyDelim)
	c.mu.RLock()
	defer c.mu.RUnlock()
	m := utils.DeepSearchInMap(c.override, paths[:len(paths)-1]...)
	dd = m[paths[len(paths)-1]]
	c.keyMap.Store(key, dd)
	return dd
}

func (c *Config) UnmarshalKey(key string, rawVal interface{}) error {
	config := mapstructure.DecoderConfig{
		DecodeHook: mapstructure.StringToTimeDurationHookFunc(),
		Result:     rawVal,
		TagName:    "toml",
	}
	decoder, err := mapstructure.NewDecoder(&config)
	if err != nil {
		return err
	}
	value := c.Get(key)
	if value == nil {
		return errors.New("invalid key, maybe not exist in config")
	}
	return decoder.Decode(value)
}

func (c *Config) traverse(sep string) map[string]interface{} {
	data := make(map[string]interface{})
	lookup("", c.override, data, sep)
	return data
}

func lookup(prefix string, target map[string]interface{}, data map[string]interface{}, sep string) {
	for k, v := range target {
		pp := fmt.Sprintf("%s%s%s", prefix, sep, k)
		if prefix == "" {
			pp = k
		}
		if dd, err := xcast.ToStringMapE(v); err == nil {
			lookup(pp, dd, data, sep)
		} else {
			data[pp] = v
		}
	}
}
