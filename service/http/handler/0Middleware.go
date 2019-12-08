package handler

import "net/http"

// Middleware ...
func Middleware(handlers ...http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, v := range handlers {
			v(w, r)
		}
	}
}
