package redis

import (
	"github.com/goinbox/pool"
)

type PConfig struct {
	pool.Config

	NewClientFunc func() (*Client, error)

	LogKeepAlive bool
}

type Pool struct {
	pl *pool.Pool

	config *PConfig
}

func NewPool(config *PConfig) *Pool {
	p := &Pool{
		config: config,
	}

	if config.NewConnFunc == nil {
		config.NewConnFunc = p.newConn
	}
	if config.KeepAliveFunc == nil {
		config.KeepAliveFunc = p.keepAlive
	}

	p.pl = pool.NewPool(&p.config.Config)

	return p
}

func (p *Pool) Get() (*Client, error) {
	conn, err := p.pl.Get()
	if err != nil {
		return nil, err
	}

	return conn.(*Client), nil
}

func (p *Pool) Put(client *Client) error {
	return p.pl.Put(client)
}

func (p *Pool) newConn() (pool.IConn, error) {
	return p.config.NewClientFunc()
}

func (p *Pool) keepAlive(conn pool.IConn) error {
	client := conn.(*Client)

	if p.config.LogKeepAlive == true {
		return client.Do("ping").Err
	}

	return client.DoWithoutLog("ping").Err
}
