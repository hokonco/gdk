package handler

import "net/http"

// GlobalOPTIONS ...
func GlobalOPTIONS() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}
}
