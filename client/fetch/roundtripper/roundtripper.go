package roundtripper

import (
	"net/http"
)

var _ http.RoundTripper = (*transport)(nil)

// New ...
func New(base http.RoundTripper, req func(r *http.Request) error, res func(r *http.Response) error) http.RoundTripper {
	return &transport{base, req, res}
}

type transport struct {
	base http.RoundTripper
	req  func(r *http.Request) error
	res  func(r *http.Response) error
}

// RoundTrip ...
func (x *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if x.base == nil {
		x.base = Transport(nil)
	}

	if x.req != nil {
		if errReq := x.req(req); errReq != nil {
			return nil, errReq
		}
	}
	res, err := x.base.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	if x.res != nil {
		if errRes := x.res(res); errRes != nil {
			return nil, errRes
		}
	}

	return res, err
}

// Transport ...
func Transport(fn func(*http.Transport)) http.RoundTripper {
	t, ok := http.DefaultTransport.(*http.Transport)
	if ok && fn != nil {
		fn(t)
	}
	return t
}
