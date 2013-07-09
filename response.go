package gadget

import (
	"fmt"
	"net/http"
)

// Response provides a wrapper around the interface{} value you would normally
// return for the response body in a Controller method, but gives you the
// ability to write headers and cookies to accompany the response.
type Response struct {
	status  int
	Body    interface{}
	final   string
	Cookies []*http.Cookie
	Headers http.Header
}

// NewResponse returns a pointer to a Response with its Body and Headers values
// initialized.
func NewResponse(body interface{}) *Response {
	return &Response{
		Body:    body,
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
	fmt.Fprint(w, r.final)
}

// AddCookie adds a cookie to the Response.
func (r *Response) AddCookie(cookie *http.Cookie) {
	r.Cookies = append(r.Cookies, cookie)
}
