package gadget

import (
	"fmt"
	"github.com/redneckbeard/gadget/env"
	"github.com/redneckbeard/quimby"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"runtime/debug"
	"strings"
	"text/tabwriter"
)

func init() {
	quimby.Add(&ListRoutes{})
}

// ListRoutes provides a command to print out all routes registered with an application.
type ListRoutes struct {
	*quimby.Flagger
}

func (c *ListRoutes) SetFlags() {}

func (c *ListRoutes) Desc() string { return "Displays list of routes registered with Gadget." }

// Run prints all the routes registered with the application.
func (c *ListRoutes) Run() {
	app.printRoutes()
}

func (a *App) printRoutes() {
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

// App provides core Gadget functionality.
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

func (a *App) GetRoutes() []*route {
	return a.routes
}

func (a *App) Host(hostname string, rtes ...*route) *route {
	rte := &route{}
	rte.hostname = regexp.MustCompile("^" + hostname + "$")
	rte.subroutes = rtes
	return rte
}

func (a *App) Mount(mountpoint string, app gdgt) *route {
	app.Configure()
	return a.Prefixed(mountpoint, app.GetRoutes()...)
}

// SetIndex creates a route that maps / to the specified controller.
// /<IdPattern> will still map to the controller's Show, Update, and Delete
// methods.
func (a *App) SetIndex(controllerName string) *route {
	route := a.newRoute(controllerName, nil)
	route.segment = ""
	route.isRoot = true
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

func (a *App) match(r *Request) (matched *route, status int, body interface{}, action string) {
	for _, route := range a.routes {
		if route.Match(r) != nil {
			matched = route
			if matched.controller == nil {
				return matched, 0, nil, ""
			}
			status, body, action := route.Respond(r)
			return matched, status, body, action
		}
	}
	return nil, 404, nil, ""
}

func (a *App) write(w http.ResponseWriter, r *Request, matched *route, status int, body interface{}, action string) {
	var response *Response
	routeData := &RouteData{
		Action: action,
		Verb:   r.Method,
	}
	if matched != nil {
		routeData.ControllerName = pluralOf(matched.controller)
	}
	contentType := r.ContentType()

	if resp, ok := body.(*Response); ok {
		response = resp
		if ct := response.Headers.Get("Content-Type"); ct != contentType && ct != "" {
			contentType = ct
		}
	} else {
		response = NewResponse(body)
	}

	status, final, mime, _ := a.process(r, status, response.Body, contentType, routeData)

	response.status = status
	response.final = final
	response.Headers.Set("Content-Type", mime)
	response.write(w)
	r.log(status, len(response.final))
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
		var final string
		req := newRequest(r)
		defer func() {
			if r := recover(); r != nil {
				trace := string(debug.Stack())
				lines := strings.Split(trace, "\n")
				trace = strings.Join(lines[6:], "\n")
				if env.Debug {
					a.write(w, req, nil, 500, trace, "")
				} else {
					a.write(w, req, nil, 500, nil, "")
					env.Log(trace)
				}
			}
		}()
		matched, status, body, action := a.match(req)
		if matched != nil && matched.handler != nil {
			matched.handler(w, r)
			return
		}
		if status == 301 || status == 302 {
			resp, ok := body.(*Response)
			if ok {
				final = resp.Body.(string)
			} else {
				final = body.(string)
				resp = NewResponse(final)
			}
			resp.Headers.Set("Location", final)
			resp.status = status
			resp.write(w)
			req.log(status, len(final))
			return
		}
		a.write(w, req, matched, status, body, action)
	}
}
