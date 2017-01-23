package webSocketHelper

import (
	"code.google.com/p/go.net/websocket"
	"errors"
	"net/http"
	"sync/atomic"
)

type Server struct {
	clientLength int32
	clients      map[int32]*websocket.Conn

	OnConnected func(index int32)
	OnClose     func()
	OnMessage   func(index int32, msg string)
}

func (this *Server) Listen(address string, port string) error {
	this.clients = make(map[int32]*websocket.Conn)
	http.Handle(address, websocket.Handler(this.onConnected))

	if "0" != port && "" != port {
		err := http.ListenAndServe(port, nil)
		return err
	}

	return nil
}

func (this *Server) onConnected(conn *websocket.Conn) {
	index := atomic.AddInt32(&this.clientLength, 1)

	defer func() {
		conn.Close()
		this.closeClient(index)
	}()

	this.clients[index] = conn

	this.OnConnected(index)

	for {
		msg := ""
		err := websocket.Message.Receive(conn, &msg)
		if err == nil {
			this.OnMessage(index, msg)
		} else {
			return
		}
	}
}

func (this *Server) closeClient(index int32) {
	delete(this.clients, index)

	this.OnClose()
}

func (this *Server) Send(client int32, msg string) error {
	if conn, bl := this.clients[client]; bl {
		_, err := conn.Write([]byte(msg))
		return err
	}
	return errors.New("client not connected")
}

func (this *Server) Boardcast(msg string) {
	b := []byte(msg)
	for _, conn := range this.clients {
		conn.Write(b)
	}
}

func (this *Server) Clients() map[int32]*websocket.Conn {
	return this.clients
}

func (this *Server) ClientLength() int {
	return len(this.clients)
}
