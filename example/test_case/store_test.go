package test_case

import (
	"context"
	"flag"
	"fmt"
	"github.com/hhzhhzhhz/gopkg/store/etcd"
	"github.com/hhzhhzhhz/gopkg/store/gorm"
	"github.com/hhzhhzhhz/gopkg/store/mysql"
	"testing"
)

func Test_store(t *testing.T) {
	t.Run("store.etcd", func(t *testing.T) {
		flag.Parse()
		c, err := etcd.RawConfig("etcd").Build()
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		resp, err := c.Get(context.TODO(), "xxx")
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		for _, v := range resp.Kvs {
			t.Log(fmt.Sprintf("%s-%s", string(v.Key), string(v.Value)))
		}
	})

	t.Run("store.mysql", func(t *testing.T) {
		flag.Parse()
		c, err := mysql.RawConfig("mysql").Build()
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		err = c.PingContext(context.TODO())
		if err != nil {
			t.Errorf(err.Error())
			return
		}
	})

	t.Run("store.gorm", func(t *testing.T) {
		flag.Parse()
		_, err := gorm.RawConfig("gorm").Build()
		if err != nil {
			t.Errorf(err.Error())
			return
		}
	})
}
