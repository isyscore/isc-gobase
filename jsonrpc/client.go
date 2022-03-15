package jsonrpc

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type Client struct {
	conn   net.Conn
	client *rpc.Client
}

func (c *Client) Listen(address string) error {
	var err error
	c.conn, err = net.Dial("tcp", address)
	if err != nil {
		return err
	}
	c.client = jsonrpc.NewClient(c.conn)
	return nil
}

func (c *Client) GetLink() *rpc.Client {
	return c.client
}

func (c *Client) Call(method string, args any, reply any) error {
	return c.client.Call(method, args, reply)
}

func (c *Client) Close() {
	_ = c.client.Close()
	_ = c.conn.Close()
}
