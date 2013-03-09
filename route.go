package gadget

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type Route struct {
	segment       string
	indexPattern  *regexp.Regexp
	objectPattern *regexp.Regexp
	actionPattern *regexp.Regexp
	controller    Controller
	subroutes     []*Route
}


func (route *Route) String() string {
	return fmt.Sprintf(route.objectPattern.String())
}

func (route *Route) buildPatterns(prefix string) {
	// Don't bother generating fancy regexps if we're looking at '/'
	if route.segment == "" {
		route.indexPattern = regexp.MustCompile(`^$`)
	} else {
		basePattern := prefix + route.segment
		route.indexPattern = regexp.MustCompile("^" + basePattern + "$")
		patternWithId := fmt.Sprintf(`^%s(?:/(?P<%s_id>%s))?$`, basePattern, route.segment, route.controller.IdPattern())
		route.objectPattern = regexp.MustCompile(patternWithId)
		actions := route.controller.ExtraActionNames()
		if len(actions) > 0 {
			actionPatternString := fmt.Sprintf(`^%s/(?:%s)$`, basePattern, strings.Join(actions, "|"))
			route.actionPattern = regexp.MustCompile(actionPatternString)
		}
	}
	// Calls to Prefixed generate routes without controllers, and the value of prefix is already all set for those
	if route.controller != nil {
		prefix += fmt.Sprintf(`%s/(?P<%s_id>%s)/`, route.segment, route.segment, route.controller.IdPattern())
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
	controller, err := Get(segment)
	if err != nil {
		panic(err)
	}
	route := &Route{segment: segment, controller: controller}
	return route
}

func (route *Route) Match(r *Request) *regexp.Regexp {
	switch {
	case route.objectPattern != nil && route.objectPattern.MatchString(r.Path):
		return route.objectPattern
	case route.actionPattern != nil && route.actionPattern.MatchString(r.Path):
		return route.actionPattern
	case route.indexPattern.MatchString(r.Path):
		return route.indexPattern
	}
	return nil
}

func (route *Route) GetParams(r *Request) map[string]string {
	params := make(map[string]string)
	pattern := route.Match(r)
	if pattern.NumSubexp() > 0 {
		names := pattern.SubexpNames()
		matches := pattern.FindStringSubmatch(r.Path)
		for i := 1; i < len(matches); i++ {
			m, n := matches[i], names[i]
			if m != "" {
				params[n] = m
			}
		}
	}
	return params
}

func (route *Route) GetActionName(r *Request) (action string) {
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

func (route *Route) Respond(r *Request) (status int, body interface{}, action string) {
	action = route.GetActionName(r)
	if action == "" {
		return 404, "", ""
	}
	r.UrlParams = route.GetParams(r)
	status, body = route.controller.RunFilters(r, action)
	if status != 0 {
		return
	}
	var methodName string
	if extra, ok := route.controller.ExtraActions()[action]; ok {
		methodName = extra
	} else {
		methodName = strings.Title(action)
	}
	t := reflect.TypeOf(route.controller)
	method, _ := t.MethodByName(methodName)
	arguments := []reflect.Value{reflect.ValueOf(route.controller), reflect.ValueOf(r)}
	statusAndBody := method.Func.Call(arguments)
	status = int(statusAndBody[0].Int())
	body = statusAndBody[1].Interface()
	return
}
