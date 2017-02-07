package webSocketHelper

import (
	"code.google.com/p/go.net/websocket"
	"gentwolf/GolangHelper/util"
	"net/http"
	"sync/atomic"
)

type Server struct {
	index        int32
	clientLength int32
	clients      map[string]*websocket.Conn

	OnConnect func(client string)
	OnClose   func(client string)
	OnMessage func(client string, b []byte)
}

//urlPath:如/chat，URL的一部分
//port:监听的端口(:80)，如果不传入，则需要在外部监听，否则使用goroutine调用此方法
func (this *Server) Listen(urlPath string, port string) error {
	this.clients = make(map[string]*websocket.Conn)
	http.Handle(urlPath, websocket.Handler(this.handleClient))

	var err error
	if "0" != port && "" != port {
		err = http.ListenAndServe(port, nil)
	}
	return err
}

func (this *Server) handleClient(conn *websocket.Conn) {
	atomic.AddInt32(&this.index, 1)
	atomic.AddInt32(&this.clientLength, 1)

	client := util.ToStr(this.index)
	this.clients[client] = conn

	defer this.CloseClient(client)

	this.OnConnect(client)

	for {
		var b []byte
		err := websocket.Message.Receive(conn, &b)
		if err == nil {
			this.OnMessage(client, b)
		} else {
			break
		}
	}
}

func (this *Server) Send(client string, b []byte) error {
	var err error
	if conn, bl := this.clients[client]; bl {
		_, err = conn.Write(b)
	}
	return err
}

func (this *Server) Boardcast(b []byte) {
	for _, conn := range this.clients {
		conn.Write(b)
	}
}

func (this *Server) ClientLength() int {
	return len(this.clients)
}

func (this *Server) CloseClient(client string) {
	if conn, bl := this.clients[client]; bl {
		conn.Close()

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
