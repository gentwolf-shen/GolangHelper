package tcpSocketHelper

import (
	"net"
	"sync/atomic"
)

type Server struct {
	clientIndex  int32
	clientLength int32
	clients      map[int32]*net.Conn

	OnConnect func(index int32)
	OnClose   func(index int32)
	OnMessage func(index int32, b []byte)
}

func (this *Server) Listen(netType, addr string) error {
	tcpAddr, err := net.ResolveTCPAddr(netType, addr)
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP(netType, tcpAddr)
	if err != nil {
		return err
	}

	this.clients = make(map[int32]*net.Conn)

	go func() {
		for {
			conn, err := listener.Accept()
			if err == nil {
				go this.handleClient(conn)
			}
		}
	}()

	return nil
}

func (this *Server) handleClient(conn net.Conn) {
	index := atomic.AddInt32(&this.clientIndex, 1)
	this.clients[index] = &conn

	atomic.AddInt32(&this.clientLength, 1)

	defer func() {
		this.CloseClient(index)
	}()

	this.OnConnect(index)

	readStream(&conn, func(b []byte) {
		this.OnMessage(index, b)
	})
}

func (this *Server) Send(index int32, b []byte) error {
	if conn, bl := this.clients[index]; bl {
		_, err := (*conn).Write(pack(b))
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *Server) CloseClient(index int32) {
	if conn, bl := this.clients[index]; bl {
		(*conn).Close()
		delete(this.clients, index)
		atomic.AddInt32(&this.clientLength, -1)
		this.OnClose(index)
	}
}

func (this *Server) CloseAll() {
	for index, _ := range this.clients {
		this.CloseClient(index)
	}
}

func (this *Server) Boardcast(b []byte) {
	for _, conn := range this.clients {
		(*conn).Write(pack(b))
	}
}

func (this *Server) Clients() map[int32]*net.Conn {
	return this.clients
}

func (this *Server) ClientLength() int {
	return int(this.clientLength)
}
