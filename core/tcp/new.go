package tcp

import (
	"dijan/core/cache"
	"dijan/core/cluster"
	"dijan/utils"
	"net"
)

type Server struct {
	cache.Cache
	cluster.Node
}

func (s *Server) Listen() {
	l, e := net.Listen("tcp", utils.GlobalSystemConfig.Server.TcpListenPort)
	if e != nil {
		panic(e)
	}
	for {
		c, e := l.Accept()
		if e != nil {
			panic(e)
		}
		go s.process(c)
	}
}

func New(c cache.Cache, n cluster.Node) *Server {
	return &Server{c, n}
}
