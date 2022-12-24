package xgin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hhzhhzhhz/gopkg/log"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
)

func NewServer(opt *Cfg) *Server {
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	return &Server{
		opt:    opt,
		Engine: gin.New(),
	}
}

type Server struct {
	opt *Cfg
	*gin.Engine
	close    int32
	Server   *http.Server
	listener net.Listener
	sync.WaitGroup
}

func (s *Server) Start() error {
	s.defaultConfig()
	addr := fmt.Sprintf(":%d", s.opt.Port)
	log.Logger().Info(fmt.Sprintf("xgin is listening and serving on %s", addr))
	s.Server = &http.Server{
		Addr:    addr,
		Handler: s,
	}
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s.listener = listener
	return s.Server.Serve(s.listener)
}

func (s *Server) Close() error {
	atomic.AddInt32(&s.close, -1)
	s.Wait()
	return nil
}

func (s *Server) defaultConfig() {
	if s.opt.StaticPath != "" {
		s.LoadHTMLGlob(s.opt.StaticPath)
	}
	if s.opt.StaticFS != "" {
		//s.r.StaticFS("/static", http.Dir("./api/static"))
		s.StaticFS("/static", http.Dir(s.opt.StaticFS))
	}
	s.NoRoute(func(c *gin.Context) {
		ResponseNo(c, Error404)
	})
	s.NoMethod(func(c *gin.Context) {
		ResponseNo(c, Error405)
	})
	s.Use(s.gracefulClose)
}

func (s *Server) gracefulClose(c *gin.Context) {
	if atomic.LoadInt32(&s.close) < 0 {
		ResponseNo(c, ErrorServerClosed)
		return
	}
	s.Add(1)
	defer func() {
		if err := recover(); err != nil {
			s.Done()
			return
		}
		s.Done()
	}()
	c.Next()
}
