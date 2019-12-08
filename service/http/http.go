package http

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hokonco/gdk/util/parallel"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/acme/autocert"
)

// Run ...
func Run(p ParameterRun) (err error) {
	var log = func(err error) {
		if err != nil {
			fmt.Println(err)
		}
	}
	if p.AddrHTTP == "" && p.AddrHTTPS == "" {
		err = fmt.Errorf("godevkit: empty http address")
		log(err)
		return
	}

	// ========================================
	// Init router
	var handler = httprouter.New()
	if p.RouterFunc != nil {
		p.RouterFunc(handler)
	}
	var srvTLS = func(isTLS bool, addr string) *http.Server {
		return &http.Server{
			IdleTimeout:  600 * time.Second,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 5 * time.Second,
			Addr:         addr,
			Handler:      handler,
			TLSConfig: func() *tls.Config {
				if !isTLS {
					return nil
				}
				return &tls.Config{
					GetCertificate: (&autocert.Manager{
						Prompt: autocert.AcceptTOS,
						Cache:  autocert.DirCache("."),
						HostPolicy: func(ctx context.Context, host string) (err error) {
							if host != p.AllowedHost {
								err = fmt.Errorf("acme/autocert: only %s host is allowed", p.AllowedHost)
							}
							return
						},
					}).GetCertificate,
				}
			}(),
		}
	}
	var srvHTTP, srvHTTPS = srvTLS(false, p.AddrHTTP), srvTLS(true, p.AddrHTTPS)

	parallel.Run(1, 0,
		func() {
			if p.AddrHTTP != "" {
				fmt.Printf("running @ http://0.0.0.0%s\n", p.AddrHTTP)
				// ========================================
				// Listen to HTTP
				err = srvHTTP.ListenAndServe()
				log(err)
			}
		},
		func() {
			if p.AddrHTTPS != "" {
				fmt.Printf("running @ https://0.0.0.0%s\n", p.AddrHTTPS)
				// ========================================
				// Listen to TLS
				err = srvHTTPS.ListenAndServeTLS("", "")
				log(err)
			}
		},
		func() {
			// ========================================
			// Listen to signal
			var chanSignal = make(chan os.Signal, 1)
			signal.Notify(chanSignal, syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM)
			err = fmt.Errorf("stopped: receiving signal [%+v]", <-chanSignal)
			log(err)
		},
	)

	// ========================================
	// Listen to shutdown
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	parallel.Run(-1, 0,
		func() {
			if p.AddrHTTP != "" {
				err = srvHTTP.Shutdown(ctx)
				log(err)
			}
		},
		func() {
			if p.AddrHTTPS != "" {
				err = srvHTTPS.Shutdown(ctx)
				log(err)
			}
		},
	)
	return
}

// ParameterRun ...
type ParameterRun struct {
	AddrHTTP    string
	AddrHTTPS   string
	AllowedHost string
	RouterFunc  func(*httprouter.Router)
}
