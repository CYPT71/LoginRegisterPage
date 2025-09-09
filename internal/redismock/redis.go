package redis

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var Nil = errors.New("redis: nil")

type Options struct {
	Addr string
}

type item struct {
	val string
	exp time.Time
}

type Client struct {
	mu   sync.RWMutex
	data map[string]item
}

func NewClient(opt *Options) *Client {
	return &Client{data: make(map[string]item)}
}

type StatusCmd struct{ err error }

func (c *StatusCmd) Err() error { return c.err }

type StringCmd struct {
	val string
	err error
}

func (c *StringCmd) Bytes() ([]byte, error)  { return []byte(c.val), c.err }
func (c *StringCmd) Result() (string, error) { return c.val, c.err }
func (c *StringCmd) Err() error              { return c.err }

type IntCmd struct{ err error }

func (c *IntCmd) Err() error { return c.err }

func (c *Client) Set(ctx context.Context, key string, value interface{}, exp time.Duration) *StatusCmd {
	c.mu.Lock()
	defer c.mu.Unlock()
	var valStr string
	switch v := value.(type) {
	case string:
		valStr = v
	case []byte:
		valStr = string(v)
	default:
		valStr = fmt.Sprint(v)
	}
	var expiry time.Time
	if exp > 0 {
		expiry = time.Now().Add(exp)
	}
	c.data[key] = item{val: valStr, exp: expiry}
	return &StatusCmd{}
}

func (c *Client) Get(ctx context.Context, key string) *StringCmd {
	c.mu.RLock()
	itm, ok := c.data[key]
	c.mu.RUnlock()
	if !ok {
		return &StringCmd{err: Nil}
	}
	if !itm.exp.IsZero() && time.Now().After(itm.exp) {
		c.mu.Lock()
		delete(c.data, key)
		c.mu.Unlock()
		return &StringCmd{err: Nil}
	}
	return &StringCmd{val: itm.val}
}

func (c *Client) Del(ctx context.Context, keys ...string) *IntCmd {
	c.mu.Lock()
	for _, k := range keys {
		delete(c.data, k)
	}
	c.mu.Unlock()
	return &IntCmd{}
}
