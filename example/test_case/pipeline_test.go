package test_case

import (
	"github.com/hhzhhzhhz/gopkg/pipeline/nsq"
	"github.com/hhzhhzhhz/gopkg/pipeline/rabbit"
	"testing"
)

func Test_MQ(t *testing.T) {
	t.Run("pipeline.nsq", func(t *testing.T) {
		c, err := nsq.RawConfig("nsq").Build()
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		if err := c.Publish("xxxx", "name"); err != nil {
			t.Errorf(err.Error())
			return
		}
	})

	t.Run("pipeline.nsq", func(t *testing.T) {
		c, err := rabbit.RawConfig("rabbit").BuildDelay()
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		if err := c.Publish("xxxx", "", 1, "xxx"); err != nil {
			t.Errorf(err.Error())
			return
		}
	})
}
