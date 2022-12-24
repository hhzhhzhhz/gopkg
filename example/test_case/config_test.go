package test_case

import (
	"github.com/hhzhhzhhz/gopkg/config"
	"testing"
)

type Db struct {
	Server        string
	Ports         []int
	ConnectionMax int
	Enabled       bool
}

func Test_Config(t *testing.T) {
	t.Log(config.GetConfig().Get("people.name"))
	db := &Db{}
	if err := config.GetConfig().UnmarshalKey("database", &db); err != nil {
		t.Errorf(err.Error())
		return
	}
	t.Log(db)
}
