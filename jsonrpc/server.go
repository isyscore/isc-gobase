package jsonrpc

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type Server struct {
	link net.Listener
}

func (s *Server) Register(v any) error {
	return rpc.Register(v)
}

func (s *Server) Listen(address string) error {
	var err error
	s.link, err = net.Listen("tcp", address)
	if err != nil {
		return err
	}
	go func() {
		for {
			if conn, err := s.link.Accept(); err == nil {
				go jsonrpc.ServeConn(conn)
			}
		}
	}()
	return nil
}

func (s *Server) Close() {
	_ = s.link.Close()
}
