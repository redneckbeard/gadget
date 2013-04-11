package gadget

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redneckbeard/gadget/env"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Request struct {
	*http.Request
	Path      string
	Params    map[string]interface{}
	Method    string
	UrlParams map[string]string
	User      User
}

func NewRequest(raw *http.Request) *Request {
	r := &Request{Request: raw, Path: raw.URL.Path[1:], Method: raw.Method}
	r.setParams()
	return r
}

func (r *Request) ContentType() string {
	accept := r.Request.Header.Get("Accept")
	if accept != "" {
		return accept
	}
	return r.Request.Header.Get("Content-Type")
}

func unpackValues(params map[string]interface{}, values map[string][]string) {
	for k, v := range values {
		if len(v) == 1 {
			params[k] = v[0]
		} else {
			params[k] = v
		}
	}
}

func (r *Request) setParams() {
	params := make(map[string]interface{})
	switch ct := r.Request.Header.Get("Content-Type"); {
	case ct == "application/json":
		if r.Request.Body != nil {
			raw, err := ioutil.ReadAll(r.Request.Body)
			defer r.Request.Body.Close()
			if err != nil {
				return
			}
			err = json.Unmarshal(raw, params)
			if err != nil {
				return
			}
		}
	case strings.HasPrefix(ct, "multipart/form-data"):
		err := r.ParseMultipartForm(10 * 1024 * 1024)
		if err != nil {
			return
		}
		unpackValues(params, r.MultipartForm.Value)
	default:
		err := r.ParseForm()
		if err != nil {
			return
		}
		unpackValues(params, r.Form)
	}
	r.Params = params
}

func (r *Request) setUser() error {
	if r.UrlParams == nil {
		return errors.New("UrlParams must be set prior to user identification")
	}
	if identifyUser != nil {
		r.User = identifyUser(r)
	} else {
		r.User = &AnonymousUser{}
	}
	return nil
}

func (r *Request) Log(status, contentLength int) {
	raw := r.Request
	env.Log(fmt.Sprintf(`[%s] "%s %s %s" %d %d`, time.Now().Format(time.RFC822), r.Method, raw.URL.Path, raw.Proto, status, contentLength))
}
