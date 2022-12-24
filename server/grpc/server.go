package grpc

import (
	"google.golang.org/grpc"
	"net"
)

// Server ... todo
type Server struct {
	*grpc.Server
	listener net.Listener
	*Cfg
}

func NewServer(opt *Cfg) *Server {
	return &Server{Cfg: opt}
}
