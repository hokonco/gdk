package web

import "net/http"

type handlerFunc struct {
	poolSize int

	method string
	path   string
	fn     http.HandlerFunc
}

// HandlerFuncs ...
type HandlerFuncs []handlerFunc

// Add new handler func
func (h HandlerFuncs) Add(poolSize int, method string, path string, fn http.HandlerFunc) HandlerFuncs {
	if fn != nil {
		h = append(h, handlerFunc{poolSize, method, path, fn})
	}
	return h
}
