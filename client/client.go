package client

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/hokonco/gdk/client/fetch"
	"github.com/hokonco/gdk/client/fetch/roundtripper"
	"github.com/hokonco/gdk/client/redis"
	"github.com/hokonco/gdk/client/sql"
)

// Fetch alias
type Fetch fetch.Config

// New Fetch
func (c Fetch) New() fetch.Instance { return fetch.New(fetch.Config(c)) }

// Redis alias
type Redis redis.Config

// New Redis
func (c Redis) New() redis.Instance { return redis.New(redis.Config(c)) }

// SQL alias
type SQL sql.Config

// New SQL
func (c SQL) New() sql.Instance { return sql.New(sql.Config(c)) }

// alias
var (
	RedisArgument = redis.Argument
	HTTPTransport = roundtripper.Transport

	HTTPDefaultTransport = roundtripper.Transport(nil)
	HTTPDumpTransport    = roundtripper.New(
		HTTPDefaultTransport,
		func(r *http.Request) error {
			bs, err := httputil.DumpRequestOut(r, true)
			fmt.Println(string(bs))
			return err
		},
		func(r *http.Response) error {
			bs, err := httputil.DumpResponse(r, true)
			fmt.Println(string(bs))
			return err
		},
	)
)
