package gadget

import (
	"fmt"
	"github.com/redneckbeard/gadget/processor"
	"net/http"
)

var routes []*Route

func Routes(rtes ...*Route) {
	for _, r := range rtes {
		routes = append(routes, r.Flatten()...)
	}
}

func SetIndex(controllerName string) *Route {
	route := newRoute(controllerName)
	route.segment = ""
	route.buildPatterns("")
	return route
}

func Resource(controllerName string, nested ...*Route) *Route {
	route := newRoute(controllerName)
	route.subroutes = nested
	route.buildPatterns("")
	return route
}

func Prefixed(prefix string, nested ...*Route) *Route {
	route := newRoute("")
	route.subroutes = nested
	route.buildPatterns(prefix)
	return route
}

func Handler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			status int
			final  string
		)
		req := NewRequest(r)
		matched := false
		for _, route := range routes {
			if route.Match(req) != nil {
				var (
					action string
					body   interface{}
				)
				matched = true
				status, body, action = route.Respond(req)
				if status == 301 || status == 302 {
					final = body.(string)
					http.Redirect(w, r, final, status)
				} else {
					var mime string
					contentType := req.ContentType()
					status, final, mime, _ = processor.Process(status, body, contentType, &processor.RouteData{Action: action, ControllerName: PluralOf(route.controller), Verb: r.Method})
					w.Header().Set("Content-Type", mime)
				}
				break
			}
		}
		if !matched {
			status = 404
			final = ""
		}
		w.WriteHeader(status)
		fmt.Fprint(w, final)
		req.Log(status, len(final))
	}
}
