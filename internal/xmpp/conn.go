package xmpp

import (
	"io"
	"net"
	"sync"
	"time"
)

const writeTimeout = 10 * time.Second

type conn struct {
	raw net.Conn
	mu  sync.Mutex
}

func (c *conn) write(text string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	_ = c.raw.SetWriteDeadline(time.Now().Add(writeTimeout))
	_, err := io.WriteString(c.raw, text)
	return err
}

func (c *conn) close() {
	c.raw.Close()
}
