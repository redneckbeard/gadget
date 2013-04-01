package gadget

import (
	"fmt"
	"net/http"
)

type Response struct {
	status  int
	Body    interface{}
	Final   string
	Cookies []*http.Cookie
	Headers http.Header
}

func NewResponse(body interface{}) *Response {
	return &Response{
		Body: body,
		Headers: make(http.Header),
	}
}

func (r *Response) write(w http.ResponseWriter) {
	h := w.Header()
	for name, values := range r.Headers {
		h[name] = values
	}
	for _, c := range r.Cookies {
		http.SetCookie(w, c)
	}
	w.WriteHeader(r.status)
	fmt.Fprint(w, r.Final)
}

func (r *Response) AddCookie(cookie *http.Cookie) {
	r.Cookies = append(r.Cookies, cookie)
}
