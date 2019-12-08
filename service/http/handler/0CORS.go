package handler

import (
	"net/http"
	"strings"
)

// CORS ...
func CORS() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var join = func(str ...string) string { return strings.Join(str, ", ") }
		for k, v := range map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": join("POST", "GET", "OPTIONS", "PUT", "DELETE"),
			"Access-Control-Allow-Headers": join("Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"),
		} {
			w.Header().Set(k, v)
		}
	}
}
