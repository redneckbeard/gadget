package routing

import (
	"fmt"
	"github.com/redneckbeard/gadget/controller"
	"github.com/redneckbeard/gadget/requests"
	"regexp"
)

type Route struct {
	segment       string
	indexPattern  *regexp.Regexp
	objectPattern *regexp.Regexp
	controller    controller.Controller
	subroutes     []*Route
}

func (route *Route) String() string {
	return fmt.Sprintf(route.objectPattern.String())
}

func (route *Route) buildPatterns(prefix string) {
	// Don't bother generating fancy regexps if we're looking at '/'
	if route.segment == "" {
		route.indexPattern = regexp.MustCompile(`^$`)
		route.objectPattern = regexp.MustCompile(`^$`)
	} else {
		basePattern := prefix + route.segment
		route.indexPattern = regexp.MustCompile("^" + basePattern + "$")
		patternWithId := fmt.Sprintf(`^%s(?:/(?P<%s_id>\d+))?$`, basePattern, route.segment)
		route.objectPattern = regexp.MustCompile(patternWithId)
	}
	// Calls to Prefixed generate routes without controllers, and the value of prefix is already all set for those
	if route.controller != nil {
		prefix += fmt.Sprintf(`%s/(?P<%s_id>\d+)/`, route.segment, route.segment)
	} else {
		prefix += "/"
	}
	for _, r := range route.subroutes {
		r.buildPatterns(prefix)
	}
}

func (route *Route) Flatten() []*Route {
	var flattened []*Route
	if route.controller != nil {
		flattened = append(flattened, route)
	}
	for _, r := range route.subroutes {
		flattened = append(flattened, r.Flatten()...)
	}
	return flattened
}

func newRoute(segment string) *Route {
	if segment == "" {
		return &Route{segment: segment}
	}
	controller, err := controller.Get(segment)
	if err != nil {
		panic(err)
	}
	route := &Route{segment: segment, controller: controller}
	return route
}

func (route *Route) Match(r *requests.Request) bool {
	if route.objectPattern.MatchString(r.Path) {
		return true
	}
	return false
}

func (route *Route) GetParams(path string) map[string]string {
	params := make(map[string]string)
	names := route.objectPattern.SubexpNames()
	matches := route.objectPattern.FindStringSubmatch(path)
	for i := 1; i < len(matches); i++ {
		m, n := matches[i], names[i]
		if m != "" {
			params[n] = m
		}
	}
	return params
}

func (route *Route) Respond(r *requests.Request) (int, interface{}) {
	r.UrlParams = route.GetParams(r.Path)
	atIndex := route.indexPattern.MatchString(r.Path)
	switch {
	case atIndex && r.Method == "GET":
		return route.controller.Index(r)
	case atIndex && r.Method == "POST":
		return route.controller.Create(r)
	case !atIndex && r.Method == "GET":
		return route.controller.Show(r)
	case !atIndex && r.Method == "PUT":
		return route.controller.Update(r)
	case !atIndex && r.Method == "DELETE":
		return route.controller.Destroy(r)
	}
	return 404, ""
}
