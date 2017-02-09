package jsonRpcHelper

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type Client struct {
	conn   net.Conn
	client *rpc.Client
}

func (this *Client) Listen(address string) error {
	var err error
	this.conn, err = net.Dial("tcp", address)
	if err != nil {
		return err
	}

	this.client = jsonrpc.NewClient(this.conn)
	return nil
}

func (this *Client) GetLink() *rpc.Client {
	return this.client
}

func (this *Client) Close() {
	this.conn.Close()
	this.client.Close()
}
