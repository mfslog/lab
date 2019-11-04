package client

import (
	"github.com/mfslog/lab/go/udpchatroom/pkg/pack"
	"github.com/mfslog/lab/go/udpchatroom/pkg/transport"
)

type Client struct {
	t *transport.Transport
}

func NewClient(t *transport.Transport) *Client {
	return &Client{t: t}
}

func (c *Client) call(pack *pack.Pack) *pack.Pack {

	return
}
