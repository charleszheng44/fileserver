package client

import (
	"github.com/sirupsen/logrus"
)

type Client struct {
	addr string
}

func NewClient(addr string) *Client {
	return &Client{
		addr: addr,
	}
}

func (c *Client) UpLoad(filePath string) error {
	return nil
}

func (c *Client) Download(filePath string) error {
	return nil
}
