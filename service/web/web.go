package web

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hokonco/gdk/pool"
	"github.com/julienschmidt/httprouter"
)

// Config ...
type Config struct {
	Addr           string
	OnShutdown     func()
	GlobalPoolSize int

	GlobalOPTIONS          http.HandlerFunc
	HandleMethodNotAllowed bool
	HandleOPTIONS          bool
	MethodNotAllowed       http.HandlerFunc
	NotFound               http.HandlerFunc
	PanicHandler           func(interface{}) http.HandlerFunc
	RedirectFixedPath      bool
	RedirectTrailingSlash  bool
	HandlerFuncs           HandlerFuncs

	ShutdownTimeout   time.Duration
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
}

// New ...
func New(conf Config, stop func() error) error {
	r := httprouter.New()
	r.GlobalOPTIONS = conf.GlobalOPTIONS
	r.HandleMethodNotAllowed = conf.HandleMethodNotAllowed
	r.HandleOPTIONS = conf.HandleOPTIONS
	r.MethodNotAllowed = conf.MethodNotAllowed
	r.NotFound = conf.NotFound
	r.PanicHandler = func(w http.ResponseWriter, r *http.Request, rcv interface{}) {
		if conf.PanicHandler != nil {
			p := conf.PanicHandler(rcv)
			if p != nil {
				p(w, r)
			}
		}
	}
	r.RedirectFixedPath = conf.RedirectFixedPath
	r.RedirectTrailingSlash = conf.RedirectTrailingSlash

	var globalPool = pool.New(conf.GlobalPoolSize)

	for _, h := range conf.HandlerFuncs {
		r.HandlerFunc(h.method, h.path, MiddlewarePool(
			h.fn,
			pool.New(h.poolSize),
			globalPool,
		))
	}
	s := &http.Server{
		Addr:              conf.Addr,
		Handler:           r,
		ReadTimeout:       conf.ReadTimeout,
		ReadHeaderTimeout: conf.ReadHeaderTimeout,
		WriteTimeout:      conf.WriteTimeout,
		IdleTimeout:       conf.IdleTimeout,
	}
	s.RegisterOnShutdown(conf.OnShutdown)

	c := make(chan error, 1)
	if stop != nil {
		go func() { c <- stop() }() // blocking to listen on stop
	}
	go func() { c <- s.ListenAndServe() }() // blocking to listen on TCP
	fmt.Printf("running @ http://localhost%+v\n", conf.Addr)
	err := <-c

	ctx, cancel := context.WithTimeout(context.Background(), conf.ShutdownTimeout)
	defer cancel()

	if err0 := s.Shutdown(ctx); err0 != nil {
		err = err0
	}
	return err
}
