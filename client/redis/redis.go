package redis

import (
	"crypto/tls"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/hokonco/gdk/errors"
	"github.com/hokonco/gdk/pool"
)

// Instance ...
type Instance interface {
	Close() error
	Do(ArgumentFunc) (interface{}, error)
	Pipeline(...ArgumentFunc) ([]interface{}, error)
	Stats() redis.PoolStats
}

type impl struct {
	*redis.Pool
	redis.Conn
	pool.Instance
}

// Config ...
type Config struct {
	Address         string
	TOBPingDuration time.Duration

	DialReadTimeout    time.Duration
	DialWriteTimeout   time.Duration
	DialConnectTimeout time.Duration
	DialKeepAlive      time.Duration
	DialTLSConfig      *tls.Config
	DialTLSSkipVerify  bool
	DialUseTLS         bool

	PoolSize            int
	PoolIdleTimeout     time.Duration
	PoolMaxIdle         int
	PoolMaxActive       int
	PoolMaxConnLifetime time.Duration
	PoolWait            bool
}

// New ...
func New(conf Config) Instance {
	var Pool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", conf.Address,
				redis.DialWriteTimeout(conf.DialReadTimeout),
				redis.DialWriteTimeout(conf.DialWriteTimeout),
				redis.DialConnectTimeout(conf.DialConnectTimeout),
				redis.DialKeepAlive(conf.DialKeepAlive),
				redis.DialTLSConfig(conf.DialTLSConfig),
				redis.DialTLSSkipVerify(conf.DialTLSSkipVerify),
				redis.DialUseTLS(conf.DialUseTLS),
			)
		},
		TestOnBorrow:    tobPing(conf.TOBPingDuration),
		IdleTimeout:     conf.PoolIdleTimeout,
		MaxIdle:         conf.PoolMaxIdle,
		MaxActive:       conf.PoolMaxActive,
		MaxConnLifetime: conf.PoolMaxConnLifetime,
		Wait:            conf.PoolWait,
	}
	if Pool == nil {
		panic(errors.New("redis: Pool[%+v]", Pool))
	}

	var Conn = Pool.Get()
	if Conn == nil {
		panic(errors.New("redis: Conn[%+v] ", Conn))
	}

	var p = pool.New(conf.PoolSize)
	if p == nil {
		panic(errors.New("redis: p[%+v]", p))
	}

	return &impl{Pool, Conn, p}
}

func (x *impl) Close() error {
	return x.Pool.Close()
}

func (x *impl) Do(fn ArgumentFunc) (interface{}, error) {
	x.Init()
	defer x.Done()

	c := x.conn()
	if fn == nil {
		fn = Argument("")
	}
	cmdName, args := fn()
	return c.Do(cmdName, args...)
}

func (x *impl) Pipeline(fns ...ArgumentFunc) ([]interface{}, error) {
	x.Init()
	defer x.Done()

	c := x.conn()
	for _, fn := range fns {
		if fn == nil {
			continue
		}
		cmdName, args := fn()
		if err := c.Send(cmdName, args...); err != nil {
			return nil, err
		}
	}
	return redis.Values(c.Do(""))
}

func (x *impl) Stats() redis.PoolStats {
	return x.Pool.Stats()
}

func (x *impl) conn() redis.Conn {
	if x.Conn.Err() != nil {
		x.Conn = x.Pool.Get()
	}
	return x.Conn
}

// ArgumentFunc pack commandName and arguments of each command
type ArgumentFunc func() (string, []interface{})

// Argument wrapper
func Argument(cmdName string, args ...interface{}) ArgumentFunc {
	return func() (string, []interface{}) { return cmdName, args }
}

func tobPing(d time.Duration) func(c redis.Conn, t time.Time) error {
	return func(c redis.Conn, t time.Time) error {
		if time.Since(t) < d {
			return nil
		}
		_, err := c.Do("PING")
		return err
	}
}
