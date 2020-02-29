package signal

import (
	"os"
	"os/signal"
)

// Listen will block until receiving signal from input
func Listen(sigs ...os.Signal) os.Signal {
	var signalChannel = make(chan os.Signal, 1)
	signal.Notify(signalChannel, sigs...)
	return <-signalChannel
}

// Wrap ...
func Wrap(sigs ...os.Signal) []os.Signal { return sigs }
