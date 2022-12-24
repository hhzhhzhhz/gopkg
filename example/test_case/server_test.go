package test_case

import (
	"flag"
	"github.com/hhzhhzhhz/gopkg/server/xgin"
	"testing"
)

func Test_service(t *testing.T) {
	flag.Parse()
	api, err := xgin.RawConfig("http").Build()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if err := api.Start(); err != nil {
		t.Errorf(err.Error())
		return
	}

}
