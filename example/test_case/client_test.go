package test_case

import (
	"flag"
	"github.com/hhzhhzhhz/gopkg/client/nsq_admin"
	"testing"
)

func Test_client(t *testing.T) {
	// need go 1.19
	//t.Start("client.kube", func(t *testing.T) {
	//  flag.Parse()
	//	_, err := kube.RawConfig("kube").Build()
	//	if err != nil {
	//		t.Errorf(err.Error())
	//		return
	//	}
	//})
	t.Run("client.nsq", func(t *testing.T) {
		flag.Parse()
		_, err := nsq_admin.RawConfig("kube").Build()
		if err != nil {
			t.Errorf(err.Error())
			return
		}
	})
}
