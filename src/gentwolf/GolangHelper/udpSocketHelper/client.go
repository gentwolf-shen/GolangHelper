package udpSocketHelper

import (
	"net"
	"time"
)

type Client struct {
	conn *net.UDPConn

	IsConnected bool
	OnMessage   func(b []byte)
	OnClose     func()
}

func (this *Client) Dial(netType, addr string) error {
	udpAddr, err := net.ResolveUDPAddr(netType, addr)
	if err != nil {
		return err
	}

	this.conn, err = net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return err
	}

	this.IsConnected = true

	go this.waitMessage()

	return nil
}

func (this *Client) waitMessage() {
	for {
		if !this.IsConnected {
			break
		}

		buf := make([]byte, BufferLength)
		n, err := this.conn.Read(buf)
		if err == nil {
			this.OnMessage(buf[0:n])
		} else {
			time.Sleep(500 * time.Microsecond)
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
