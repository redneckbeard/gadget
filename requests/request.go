package requests

import (
	"fmt"
	"net/http"
	"github.com/redneckbeard/gadget/env"
	"time"
)

type Request struct {
	*http.Request
	Path      string
	Method    string
	UrlParams map[string]string
}

func New(raw *http.Request) *Request {
	return &Request{Request: raw, Path: raw.URL.Path[1:], Method: raw.Method}
}

func (r *Request) ContentType() string {
	accept := r.Request.Header.Get("Accept")
	if accept != "" {
		return accept
	}
	return r.Request.Header.Get("Content-Type")
}

func (r *Request) Log(status, contentLength int) {
	raw := r.Request
	env.Log(fmt.Sprintf(`[%s] "%s %s %s" %d %d`, time.Now().Format(time.RFC822), r.Method, raw.URL.Path, raw.Proto, status, contentLength))
}
