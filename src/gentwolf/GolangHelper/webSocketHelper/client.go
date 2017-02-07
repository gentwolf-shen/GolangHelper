package webSocketHelper

import (
	"code.google.com/p/go.net/websocket"
	"strings"
)

type Client struct {
	addr string
	conn *websocket.Conn

	IsConnected bool
	OnMessage   func(b []byte)
	OnClose     func()
}

//addr:如 ws://127.0.0.1:9997/chat
func (this *Client) Dial(addr string) error {
	this.addr = addr

	err := this.ReDial()
	if err == nil {
		go this.waitMessage()
	}

	return err
}

func (this *Client) ReDial() error {
	var err error
	this.conn, err = websocket.Dial(this.addr, "", strings.Replace(this.addr, "ws", "http", -1))
	if err == nil {
		this.IsConnected = true
	}
	return err
}

func (this *Client) waitMessage() {
	defer this.Close()

	for {
		if this.IsConnected {
			var b []byte
			err := websocket.Message.Receive(this.conn, &b)
			if err == nil {
				this.OnMessage(b)
			}
		}
	}
}

func (this *Client) Send(b []byte) error {
	_, err := this.conn.Write(b)
	return err
}

func (this *Client) Close() {
	this.IsConnected = false
	this.conn.Close()
	this.OnClose()
}
