package mysql

import (
	"github.com/goinbox/pool"
)

type PConfig struct {
	pool.Config

	NewClientFunc func() (*Client, error)
}

type Pool struct {
	pl *pool.Pool

	config *PConfig
}

type NewClientFunc func() (*Client, error)

func NewPool(config *PConfig) *Pool {
	p := &Pool{
		config: config,
	}

	if config.NewConnFunc == nil {
		config.NewConnFunc = p.newConn
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
