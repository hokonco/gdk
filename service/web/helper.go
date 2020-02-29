package web

import (
	"context"
	"net/http"

	"github.com/hokonco/gdk/pool"
	"github.com/julienschmidt/httprouter"
)

// MiddlewarePool ...
func MiddlewarePool(next http.HandlerFunc, p pool.Instance, g pool.Instance) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pool.Init(p, g)
		defer pool.Done(p, g)
		next(w, r)
	}
}

// ParamsFromContext ...
func ParamsFromContext(ctx context.Context) httprouter.Params {
	return httprouter.ParamsFromContext(ctx)
}
