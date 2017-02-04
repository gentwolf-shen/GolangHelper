package udpSocketHelper

import (
	"net"
	"sync/atomic"
)

type Server struct {
	clientLength int32
	clients      map[string]*net.UDPAddr

	conn *net.UDPConn

	IsConnected bool
	OnConnect   func(client string)
	OnClose     func(client string)
	OnMessage   func(client string, b []byte)
}

func (this *Server) Listen(netType, addr string) error {
	udpAddr, err := net.ResolveUDPAddr(netType, addr)
	if err != nil {
		return err
	}

	this.conn, err = net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}

	this.clients = make(map[string]*net.UDPAddr)
	this.IsConnected = true

	go this.handleClient()

	return nil
}

//接收客户端消息，并保存客户端信息
func (this *Server) handleClient() {
	for {
		if !this.IsConnected {
			break
		}

		buf := make([]byte, BufferLength)
		if n, clientAddr, err := this.conn.ReadFromUDP(buf); err == nil {
			client := clientAddr.String()
			_, bl := this.clients[client]
			if !bl {
				atomic.AddInt32(&this.clientLength, 1)
				this.clients[client] = clientAddr
				this.OnConnect(client)
			}
			this.OnMessage(client, buf[0:n])
		}
	}
}

//发送消息到客户端，如果发送失败，则将客户端信息删除
func (this *Server) Send(client string, b []byte) error {
	if addr, bl := this.clients[client]; bl {
		_, err := this.conn.WriteToUDP(b, addr)
		if err != nil {
			delete(this.clients, client)
			atomic.AddInt32(&this.clientLength, -1)
			this.OnClose(client)
		}
		return nil
	}
	return nil
}

func (this *Server) Boardcast(b []byte) {
	for client, _ := range this.clients {
		this.Send(client, b)
	}
}

func (this *Server) Close() {
	this.IsConnected = false
	this.conn.Close()
	this.clientLength = 0
	this.clients = nil
}
