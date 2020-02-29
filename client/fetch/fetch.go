package fetch

import (
	"net/http"
	"time"

	"github.com/hokonco/gdk/pool"
)

// Instance ...
type Instance interface {
	Do(*http.Request) (*http.Response, error)
}

type impl struct {
	*http.Client
	pool.Instance
}

// Config ...
type Config struct {
	Transport http.RoundTripper
	Timeout   time.Duration

	PoolSize int
}

// New ...
func New(conf Config) Instance {
	c := &http.Client{
		Transport: conf.Transport,
		Timeout:   conf.Timeout,
	}

	return &impl{c, pool.New(conf.PoolSize)}
}

func (x *impl) Do(req *http.Request) (*http.Response, error) {
	x.Init()
	defer x.Done()
	defer x.Client.CloseIdleConnections()

	return x.Client.Do(req)
}
