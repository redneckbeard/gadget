package routing

import (
	"fmt"
	"github.com/redneckbeard/gadget/controller"
	"github.com/redneckbeard/gadget/requests"
	"reflect"
	"regexp"
	"strings"
)

type Route struct {
	segment       string
	indexPattern  *regexp.Regexp
	objectPattern *regexp.Regexp
	actionPattern *regexp.Regexp
	controller    controller.Controller
	subroutes     []*Route
}

func arbitraryActions(controller controller.Controller) (actions []string) {
	t := reflect.TypeOf(controller)
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		if method.PkgPath == "" {
			actions = append(actions, strings.ToLower(method.Name))
		}
	}
	return
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
		actions := arbitraryActions(route.controller)
		if len(actions) > 0 {
			actionPatternString := fmt.Sprintf(`^%s/(?:%s)$`, basePattern, strings.Join(actions, "|"))
			route.actionPattern = regexp.MustCompile(actionPatternString)
		}
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

func (route *Route) GetActionName(r *requests.Request) (action string) {
	atIndex := route.indexPattern.MatchString(r.Path)
	switch {
	case route.actionPattern != nil && route.actionPattern.MatchString(r.Path):
		segments := strings.Split(r.Path, "/")
		action = segments[len(segments)-1]
	case atIndex && r.Method == "GET":
		action = "index"
	case atIndex && r.Method == "POST":
		action = "create"
	case !atIndex && r.Method == "GET":
		action = "show"
	case !atIndex && r.Method == "PUT":
		action = "update"
	case !atIndex && r.Method == "DELETE":
		action = "destroy"
	}
	return
}

func (route *Route) Respond(r *requests.Request) (status int, body interface{}, action string) {
	action = route.GetActionName(r)
	if action == "" {
		return 404, "", ""
	}
	r.UrlParams = route.GetParams(r.Path)
	t := reflect.TypeOf(route.controller)
	methodName := strings.Title(action)
	method, _ := t.MethodByName(methodName)
	arguments := []reflect.Value{reflect.ValueOf(route.controller), reflect.ValueOf(r)}
	statusAndBody := method.Func.Call(arguments)
	status = int(statusAndBody[0].Int())
	body = statusAndBody[1].Interface()
	return
}
