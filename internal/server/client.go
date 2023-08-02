package server

import (
	"context"
	"io"
	"strings"
	"time"
)

type client struct {
	rw   io.ReadWriter
	msgs chan string
	err  chan error
}

func NewClient(rw io.ReadWriter) *client {
	c := &client{
		rw:   rw,
		msgs: make(chan string),
		err:  make(chan error),
	}

	go c.run()
	return c
}

func (c *client) ReadWithTimeout(ctx context.Context, timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.Read(ctx)
}

func (c *client) Read(ctx context.Context) (string, error) {
	select {
	case msg := <-c.msgs:
		return msg, nil
	case err := <-c.err:
		return "", err
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func (c *client) Write(msg string) error {
	_, err := c.rw.Write([]byte(msg))
	return err
}

func (c *client) run() {
	defer close(c.msgs)
	defer close(c.err)

	buf := make([]byte, 1024)

	for {
		n, err := c.rw.Read(buf)
		if err != nil {
			c.err <- err
			return
		}

		c.msgs <- strings.TrimSpace(string(buf[:n]))
	}
}
