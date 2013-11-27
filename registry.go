package gadget

import (
	"fmt"
	"github.com/redneckbeard/quimby"
	"net/http"
	"os"
	"reflect"
	"text/tabwriter"
)

func init() {
	quimby.Add(&ListRoutes{})
}

type ListRoutes struct {
	*quimby.Flagger
}

func (c *ListRoutes) SetFlags() {}

func (c *ListRoutes) Desc() string { return "Displays list of routes registered with Gadget." }

func (c *ListRoutes) Run() {
	app.PrintRoutes()
}

func (a *App) PrintRoutes() {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)

	for _, r := range a.routes {
		var desc string
		if r.controller != nil {
			t := reflect.TypeOf(r.controller)
			desc = t.String()[1:]
		} else {
			desc = "http.HandlerFunc"
		}
		fmt.Fprintln(w, r, "\t", desc)
	}
	w.Flush()
}

type App struct {
	routes      []*route
	Brokers     map[string]Broker
	Controllers map[string]Controller
}

// Routes registers a variable number of routes with the Gadget router. Arguments to
// Routes should be calls to SetIndex, Resource, or Prefixed.
func (a *App) Routes(rtes ...*route) {
	a.routes = []*route{}
	for _, r := range rtes {
		a.routes = append(a.routes, r.flatten()...)
	}
}

// SetIndex creates a route that maps / to the specified controller.
// /<IdPattern> will still map to the controller's Show, Update, and Delete
// methods.
func (a *App) SetIndex(controllerName string) *route {
	route := a.newRoute(controllerName, nil)
	route.segment = ""
	route.buildPatterns("")
	return route
}

// Resource creates a route to the specified controller and optionally creates
// additional routes nested under it.
func (a *App) Resource(controllerName string, nested ...*route) *route {
	route := a.newRoute(controllerName, nil)
	route.subroutes = nested
	route.buildPatterns("")
	return route
}

// Prefixed mounts routes at a URL path that is not necessarily a controller
// name.
func (a *App) Prefixed(prefix string, nested ...*route) *route {
	route := a.newRoute("", nil)
	route.subroutes = nested
	route.buildPatterns(prefix)
	return route
}

// HandleFunc mounts an http.HandlerFunc at the specified URL.
func (a *App) HandleFunc(mount string, handler http.HandlerFunc) *route {
	route := a.newRoute(mount, handler)
	route.buildPatterns("")
	return route
}

// Handler returns a func encapsulating the Gadget router (and corresponding
// controllers that can be used in a call to http.HandleFunc. Handler must be
// invoked only after Routes has been called and all Controllers have been
// registered.
//
// In theory, Gadget users will not ever have to call Handler, as Gadget will
// set up http.HandleFunc to use its return value.
func (a *App) Handler() http.HandlerFunc {
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
		for _, route := range a.routes {
			if route.Match(req) != nil {
				if route.handler != nil {
					route.handler(w, r)
					return
				}
				matched = route
				status, body, action = route.Respond(req)
				if status == 301 || status == 302 {
					resp, ok := body.(*Response)
					if ok {
						final = resp.Body.(string)
					} else {
						final = body.(string)
					}
					resp.Headers.Set("Location", final)
					resp.status = status
					resp.write(w)
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
			if ct := response.Headers.Get("Content-Type"); ct != contentType && ct != "" {
				contentType = ct
			}
		} else {
			response = NewResponse(body)
		}

		status, final, mime, _ := a.Process(req, status, response.Body, contentType, routeData)

		response.status = status
		response.final = final
		response.Headers.Set("Content-Type", mime)
		response.write(w)
		req.log(status, len(response.final))
	}
}
