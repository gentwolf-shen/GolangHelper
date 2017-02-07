package tcpSocketHelper

import (
	"net"
	"sync/atomic"
)

type Server struct {
	clientLength int32
	clients      map[string]*net.Conn

	OnConnect func(client string)
	OnClose   func(client string)
	OnMessage func(client string, b []byte)
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

	this.clients = make(map[string]*net.Conn)

	go func() {
		for {
			conn, err := listener.Accept()
			if err == nil {
				go this.handleClient(&conn)
			}
		}
	}()

	return nil
}

func (this *Server) handleClient(conn *net.Conn) {
	atomic.AddInt32(&this.clientLength, 1)

	client := (*conn).RemoteAddr().String()
	this.clients[client] = conn

	defer this.CloseClient(client)

	this.OnConnect(client)

	readStream(conn, func(b []byte) {
		this.OnMessage(client, b)
	})
}

func (this *Server) Send(client string, b []byte) error {
	var err error
	if conn, bl := this.clients[client]; bl {
		_, err = (*conn).Write(pack(b))
	}
	return err
}

func (this *Server) CloseClient(client string) {
	if conn, bl := this.clients[client]; bl {
		(*conn).Close()

		delete(this.clients, client)
		atomic.AddInt32(&this.clientLength, -1)
		this.OnClose(client)
	}
}

func (this *Server) CloseAll() {
	for client, _ := range this.clients {
		this.CloseClient(client)
	}
}

func (this *Server) Boardcast(b []byte) {
	for _, conn := range this.clients {
		(*conn).Write(pack(b))
	}
}

func (this *Server) ClientLength() int {
	return int(this.clientLength)
}
