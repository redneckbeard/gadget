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
			status  int
			final   string
			action  string
			body    interface{}
			matched *Route
		)
		req := NewRequest(r)
		for _, route := range routes {
			if route.Match(req) != nil {
				matched = route
				status, body, action = route.Respond(req)
				if status == 301 || status == 302 {
					final = body.(string)
					http.Redirect(w, r, final, status)
				}
				break
			}
		}
		routeData := &processor.RouteData{
			Action: action,
			Verb:   r.Method,
		}
		if matched == nil {
			status = 404
			final = ""
		} else {
			routeData.ControllerName = PluralOf(matched.controller)
		}
		contentType := req.ContentType()

		status, final, mime, _ := processor.Process(status, body, contentType, routeData)
		w.Header().Set("Content-Type", mime)
		w.WriteHeader(status)
		fmt.Fprint(w, final)
		req.Log(status, len(final))
	}
}
