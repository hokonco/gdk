package http

import (
	"net/http"
	"net/url"

	"github.com/julienschmidt/httprouter"
)

// Headers ...
type Headers http.Header

// Queries ...
type Queries url.Values

// Params ...
type Params httprouter.Params

// ParseHTTPRequest ...
func ParseHTTPRequest(r *http.Request) (Headers, Queries, Params) {
	return Headers(r.Header), Queries(r.URL.Query()), Params(httprouter.ParamsFromContext(r.Context()))
}
