package client

import (
	"net"
	"strings"

	"github.com/pkg/errors"
)

type Client struct {
	server string
	cli    net.Conn
}

func New(server string) *Client {
	return &Client{
		server: server,
	}
}

func (c *Client) Connect() error {
	var err error

	c.cli, err = net.Dial("tcp", c.server)
	if err != nil {
		return errors.Wrap(err, "failed to connect to server")
	}

	return nil
}

func (c *Client) Close() error {
	return c.cli.Close()
}

func (c *Client) Receive() (string, error) {
	buf := make([]byte, 1024)

	n, err := c.cli.Read(buf)
	if err != nil {
		return "", errors.Wrap(err, "failed to receive message")
	}

	return strings.TrimSpace(string(buf[:n])), nil
}

func (c *Client) Send(msg string) error {
	_, err := c.cli.Write([]byte(msg))
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}
