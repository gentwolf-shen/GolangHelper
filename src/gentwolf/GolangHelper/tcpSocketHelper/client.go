package tcpSocketHelper

import (
	"net"
)

type Client struct {
	tcpAddr *net.TCPAddr
	conn    net.Conn

	IsConnected bool
	OnMessage   func(b []byte)
	OnClose     func()
}

func (this *Client) Dial(netType, addr string) error {
	tcpAddr, err := net.ResolveTCPAddr(netType, addr)
	if err != nil {
		return err
	}

	this.tcpAddr = tcpAddr
	if err = this.ReDial(); err != nil {
		return err
	}

	go this.waitMessage()

	return nil
}

func (this *Client) ReDial() error {
	var err error
	if !this.IsConnected {
		this.conn, err = net.DialTCP("tcp", nil, this.tcpAddr)
		if err == nil {
			this.IsConnected = true
		}
	}
	return err
}

func (this *Client) waitMessage() {
	readStream(&this.conn, func(b []byte) {
		this.OnMessage(b)
	})
}

func (this *Client) Send(b []byte) error {
	_, err := this.conn.Write(pack(b))
	return err
}

func (this *Client) Close() {
	this.IsConnected = false
	this.conn.Close()
	this.OnClose()
}
