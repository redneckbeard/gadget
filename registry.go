package gadget

import (
	"fmt"
	"github.com/redneckbeard/gadget/cmd"
	"net/http"
	"os"
	"reflect"
	"text/tabwriter"
)

func init() {
	cmd.Add(&ListRoutes{})
}

var routes []*route

type ListRoutes struct {
	*cmd.Flagger
}

func (c *ListRoutes) SetFlags() {}

func (c *ListRoutes) Desc() string { return "Displays list of routes registered with Gadget." }

func (c *ListRoutes) Run() {
	PrintRoutes()
}

func PrintRoutes() {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)

	for _, r := range routes {
		t := reflect.TypeOf(r.controller)
		fmt.Fprintln(w, r, "\t", t.String()[1:])
	}
	w.Flush()
}

// Routes registers a variable number of routes with the Gadget router. Arguments to
// Routes should be calls to SetIndex, Resource, or Prefixed.
func Routes(rtes ...*route) {
	routes = []*route{}
	for _, r := range rtes {
		routes = append(routes, r.flatten()...)
	}
}

// SetIndex creates a route that maps / to the specified controller.
// /<IdPattern> will still map to the controller's Show, Update, and Delete
// methods.
func SetIndex(controllerName string) *route {
	route := newRoute(controllerName)
	route.segment = ""
	route.buildPatterns("")
	return route
}

// Resource creates a route to the specified controller and optionally creates
// additional routes nested under it.
func Resource(controllerName string, nested ...*route) *route {
	route := newRoute(controllerName)
	route.subroutes = nested
	route.buildPatterns("")
	return route
}

// Prefixed mounts routes at a URL path that is not necessarily a controller
// name.
func Prefixed(prefix string, nested ...*route) *route {
	route := newRoute("")
	route.subroutes = nested
	route.buildPatterns(prefix)
	return route
}

// Handler returns a func encapsulating the Gadget router (and corresponding
// controllers that can be used in a call to http.HandleFunc. Handler must be
// invoked only after Routes has been called and all Controllers have been
// registered.
//
// In theory, Gadget users will not ever have to call Handler, as Gadget will
// set up http.HandleFunc to use its return value.
func Handler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			status   int
			final    string
			action   string
			body     interface{}
			matched  *route
			response *Response
		)
		req := newRequest(r)
		for _, route := range routes {
			if route.Match(req) != nil {
				matched = route
				status, body, action = route.Respond(req)
				if status == 301 || status == 302 {
					if resp, ok := body.(*Response); ok {
						final = resp.Body.(string)
					} else {
						final = body.(string)
					}
					http.Redirect(w, r, final, status)
					req.log(status, len(final))
					return
				}
				break
			}
		}
		routeData := &RouteData{
			Action: action,
			Verb:   r.Method,
		}
		if matched == nil {
			status = 404
			final = ""
		} else {
			routeData.ControllerName = pluralOf(matched.controller)
		}
		contentType := req.ContentType()

		if resp, ok := body.(*Response); ok {
			response = resp
			if ct := response.Headers.Get("Content-Type"); ct != contentType {
				contentType = ct
			}
		} else {
			response = NewResponse(body)
		}

		status, final, mime, _ := Process(req, status, response.Body, contentType, routeData)

		response.status = status
		response.final = final
		response.Headers.Set("Content-Type", mime)
		response.write(w)
		req.log(status, len(response.final))
	}
}
