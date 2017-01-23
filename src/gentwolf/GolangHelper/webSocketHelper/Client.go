package webSocketHelper

import (
	"code.google.com/p/go.net/websocket"
	"strings"
)

type Client struct {
	conn *websocket.Conn

	OnMessage func(msg string)
	OnClose   func()
}

func (this *Client) Dial(address string) error {
	var err error
	this.conn, err = websocket.Dial(address, "", strings.Replace(address, "ws", "http", 0))
	if err == nil {
		go this.onMessage()
	}

	return err
}

func (this *Client) onMessage() {
	defer this.Close()

	for {
		msg := ""
		err := websocket.Message.Receive(this.conn, &msg)
		if err == nil {
			this.OnMessage(msg)
		} else {
			return
		}
	}
}

func (this *Client) Send(msg string) error {
	_, err := this.conn.Write([]byte(msg))
	return err
}

func (this *Client) Close() {
	this.conn.Close()
	this.OnClose()
}
