package gadget

import (
	"encoding/json"
	"fmt"
	"github.com/redneckbeard/gadget/env"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type (
	debugChecker  func(*Request) bool
	requestLogger func(*Request, int, int) string
)

// SetDebugWith by default is a function that always returns false, no matter
// what request is passed to it.
var (
	SetDebugWith  = func(r *Request) bool { return false }
	RequestLogger = func(r *Request, status, contentLength int) string {
		return fmt.Sprintf(`[%s] "%s %s %s" %d %d`, time.Now().Format(time.RFC822), r.Method, r.URL.Path, r.Proto, status, contentLength)
	}
)

// Request wraps an *http.Request and adds some Gadget-derived conveniences. The
// Params map contains either POST data, GET query parameters, or the body of the
// request deserialized as JSON if the request sends an Accept header of
// application/json. The UrlParams map contains any resource ids plucked from the
// URL by the router. The User is either an AnonymousUser or an object returned by
// the UserIdentifier that the application as registered with IdentifyUsersWith.
type Request struct {
	*http.Request
	Params    map[string]interface{}
	Path      string
	UrlParams map[string]string
	User      User
	RawJson   []byte
}

func newRequest(raw *http.Request) *Request {
	r := &Request{Request: raw, Path: raw.URL.Path[1:]}
	r.setParams()
	return r
}

// ContentType is sort of a dishonest method -- it returns the value of an
// Accept header if present, and falls back to Content-Type.
func (r *Request) ContentType() string {
	accept := r.Request.Header.Get("Accept")
	if accept != "" {
		return accept
	}
	return r.contentType()
}

// Debug returns true if env.Debug is true or if SetDebugWith returns true when
// passed its receiver r.
func (r *Request) Debug() bool {
	return env.Debug || SetDebugWith(r)
}

func (r *Request) Unmarshal(i interface{}) error {
	return json.Unmarshal(r.RawJson, i)
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

func (r *Request) contentType() string {
	return strings.Split(r.Request.Header.Get("Content-Type"), ";")[0]
}

func (r *Request) setParams() {
	params := make(map[string]interface{})
	switch ct := r.contentType(); {
	case ct == "application/json":
		if r.Body != nil {
			raw, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			if err != nil {
				return
			}
			err = json.Unmarshal(raw, &params)
			r.RawJson = raw
			if err != nil {
				env.Log("Unable to deserialize JSON payload: ", err, string(raw))
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
	if identifyUser != nil {
		r.User = identifyUser(r)
	} else {
		r.User = &AnonymousUser{}
	}
	return nil
}

func (r *Request) log(status, contentLength int) {
	env.Log(RequestLogger(r, status, contentLength))
}
