package service

import (
	"os"

	"github.com/hokonco/gdk/errors"
	"github.com/hokonco/gdk/service/web"
	"github.com/hokonco/gdk/signal"
)

// Config ...
type Config struct {
	ListenSignals []os.Signal
	Web           web.Config
}

// New ...
func New(conf Config) error {
	var doneChannel = make(chan struct{}, 1)
	var webErrorChannel = make(chan error, 1)
	var err error

	go func() {
		// blocking
		var sig = signal.Listen(conf.ListenSignals...)
		webErrorChannel <- errors.New("signal: %+v", sig)
	}()
	go func() {
		// blocking
		var stop = func() error { return <-webErrorChannel }
		err = web.New(conf.Web, stop)
		doneChannel <- struct{}{}
	}()

	<-doneChannel

	return err
}
