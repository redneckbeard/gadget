package requests

import "net/http"

type Request struct {
	*http.Request
	Path      string
	Method    string
	UrlParams map[string]string
}

func New(raw *http.Request) *Request {
	return &Request{Request: raw, Path: raw.URL.Path[1:], Method: raw.Method}
}
