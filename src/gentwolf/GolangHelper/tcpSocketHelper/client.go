package tcpSocketHelper

import (
	"net"
)

type Client struct {
	conn net.Conn

	OnMessage func(b []byte)
	OnClose   func()
}

func (this *Client) Dial(netType, addr string) error {
	tcpAddr, err := net.ResolveTCPAddr(netType, addr)
	if err != nil {
		return err
	}

	this.conn, err = net.DialTCP(netType, nil, tcpAddr)
	if err != nil {
		return err
	}

	go this.waitMessage()

	return nil
}

func (this *Client) waitMessage() {
	defer func() {
		this.Close()
	}()

	readStream(&this.conn, func(b []byte) {
		this.OnMessage(b)
	})
}

func (this *Client) Send(b []byte) (int, error) {
	return this.conn.Write(pack(b))
}

func (this *Client) Close() {
	this.conn.Close()
	this.OnClose()
}
