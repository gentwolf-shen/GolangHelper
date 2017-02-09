package jsonRpcHelper

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type Server struct {
	link net.Listener
}

func (this *Server) Register(v interface{}) error {
	return rpc.Register(v)
}

func (this *Server) Listen(address string) error {
	var err error

	this.link, err = net.Listen("tcp", address)
	if err != nil {
		return err
	}

	go func() {
		for {
			if conn, err := this.link.Accept(); err == nil {
				go jsonrpc.ServeConn(conn)
			}
		}
	}()

	return nil
}

func (this *Server) Close() {
	this.link.Close()
}
