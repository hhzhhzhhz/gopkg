package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hhzhhzhhz/gopkg/log"
	"github.com/hhzhhzhhz/gopkg/runtime"
	"github.com/hhzhhzhhz/gopkg/server"
	"github.com/hhzhhzhhz/gopkg/server/xgin"
	"github.com/hhzhhzhhz/gopkg/store/mysql"
	"github.com/hhzhhzhhz/gopkg/utils"
	"go.uber.org/multierr"
	"os"
)

func main() {
	if err := server.Run(&program{}); err != nil {
		fmt.Fprint(os.Stderr, fmt.Sprintf("main.run failed cause: %s", err.Error()))
		os.Exit(1)
	}
}

type Schedule struct {
	api *xgin.Server
}

func (s *Schedule) Start() error {
	return nil
}

func (s *Schedule) Close() error {
	return nil
}

type program struct {
	ctx      context.Context
	cancel   context.CancelFunc
	schedule *Schedule
	http     *xgin.Server
	utils.WaitGroupWrapper
}

func (p *program) Init() error {
	p.ctx, p.cancel = context.WithCancel(context.Background())
	log.Logger().Info("server initializing")
	flag.Parse()
	var err error
	p.http, err = xgin.RawConfig("http").Build()
	if err != nil {
		return err
	}
	p.http.GET("/", func(c *gin.Context) {
		xgin.ResponseOk(c)
	})
	p.schedule = &Schedule{api: p.http}
	_, err = mysql.RawConfig("mysql").Build()
	if err != nil {
		return err
	}
	log.Logger().Info("server initialization success")
	return nil
}

func (p *program) Start() error {
	log.Logger().Info("server starting")
	runtime.StartPprof(runtime.Pprof)
	runtime.StartMetric(runtime.Metric)
	if err := p.http.Start(); err != nil {
		return err
	}
	if err := p.schedule.Start(); err != nil {
		return err
	}
	log.Logger().Info("server started")
	return nil
}

func (p *program) Stop() error {
	var err error
	log.Logger().Info("server ready to close")
	p.cancel()
	p.Wait()
	err = multierr.Append(err, p.schedule.Close())
	err = multierr.Append(err, p.http.Close())
	log.Logger().Info("server is closed")
	err = multierr.Append(err, log.Logger().Close())
	return err
}
