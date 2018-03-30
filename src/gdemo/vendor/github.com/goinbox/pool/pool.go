package pool

import (
	"errors"
	"time"
)

type IConn interface {
	Free()
}

type Config struct {
	Size              int
	MaxIdleTime       time.Duration
	KeepAliveInterval time.Duration

	NewConnFunc   func() (IConn, error)
	KeepAliveFunc func(conn IConn) error
}

type poolItem struct {
	conn IConn

	accessTime time.Time
}

var ErrPoolIsFull = errors.New("pool is full")

type Pool struct {
	config *Config

	conns chan *poolItem
}

func NewPool(config *Config) *Pool {
	p := &Pool{
		config: config,

		conns: make(chan *poolItem, config.Size),
	}

	if config.KeepAliveInterval > 0 && config.KeepAliveFunc != nil {
		go p.keepAliveRoutine()
	}

	return p
}

func (p *Pool) Get() (IConn, error) {
	pi := p.get()
	if pi != nil {
		if time.Now().Sub(pi.accessTime) < p.config.MaxIdleTime {
			return pi.conn, nil
		}

		pi.conn.Free()
	}

	return p.config.NewConnFunc()
}

func (p *Pool) Put(conn IConn) error {
	pi := &poolItem{
		conn: conn,

		accessTime: time.Now(),
	}

	notFull := p.put(pi)
	if notFull {
		return nil
	}

	conn.Free()

	return ErrPoolIsFull
}

func (p *Pool) get() *poolItem {
	select {
	case pi := <-p.conns:
		return pi
	default:
	}

	return nil
}

func (p *Pool) put(pi *poolItem) bool {
	select {
	case p.conns <- pi:
		return true
	default:
	}

	return false
}

func (p *Pool) keepAliveRoutine() {
	ticker := time.NewTicker(p.config.KeepAliveInterval)

	for {
		select {
		case <-ticker.C:
			p.keepAlive()
		}
	}
}

func (p *Pool) keepAlive() {
	maxIdleNum := len(p.conns)

	for i := 0; i < maxIdleNum; i++ {
		pi := p.get()
		if pi == nil {
			return
		}

		if time.Now().Sub(pi.accessTime) < p.config.MaxIdleTime {
			err := p.config.KeepAliveFunc(pi.conn)
			if err == nil {
				if p.put(pi) {
					continue
				}
			}
		}

		pi.conn.Free()
	}
}
