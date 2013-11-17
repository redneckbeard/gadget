package gadget

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

type route struct {
	segment                                    string
	indexPattern, objectPattern, actionPattern *regexp.Regexp
	handler                                    http.HandlerFunc
	controller                                 Controller
	subroutes                                  []*route
}

func (rte *route) String() string {
	if rte.objectPattern != nil {
		return rte.objectPattern.String()
	}
	return rte.indexPattern.String()
}

func (rte *route) buildPatterns(prefix string) {
	// Don't bother generating fancy regexps if we're looking at '/'
	if rte.segment == "" {
		rte.indexPattern = regexp.MustCompile(`^$`)
	} else {
		basePattern := prefix + rte.segment
		rte.indexPattern = regexp.MustCompile("^" + basePattern + "$")
		if rte.controller != nil {
			patternWithId := fmt.Sprintf(`^%s(?:/(?P<%s_id>%s))?$`, basePattern, strings.Replace(nameOf(rte.controller), "-", "_", -1), rte.controller.IdPattern())
			rte.objectPattern = regexp.MustCompile(patternWithId)
			actions := rte.controller.extraActionNames()
			if len(actions) > 0 {
				actionPatternString := fmt.Sprintf(`^%s/(?:%s)$`, basePattern, strings.Join(actions, "|"))
				rte.actionPattern = regexp.MustCompile(actionPatternString)
			}
		}
	}
	// Calls to Prefixed generate routes without controllers, and the value of prefix is already all set for those
	if rte.controller != nil {
		prefix += fmt.Sprintf(`%s/(?P<%s_id>%s)/`, rte.segment, nameOf(rte.controller), rte.controller.IdPattern())
	} else {
		prefix += "/"
	}
	for _, r := range rte.subroutes {
		r.buildPatterns(prefix)
	}
}

func (rte *route) flatten() []*route {
	var flattened []*route
	if rte.controller != nil || rte.handler != nil {
		flattened = append(flattened, rte)
	}
	for _, r := range rte.subroutes {
		flattened = append(flattened, r.flatten()...)
	}
	return flattened
}

func newRoute(segment string, handler http.HandlerFunc) *route {
	if segment == "" {
		return &route{segment: segment}
	}
	controller, err := getController(segment)
	if err != nil && handler == nil {
		panic(err)
	}
	rte := &route{segment: segment, controller: controller, handler: handler}
	return rte
}

func (rte *route) Match(r *Request) *regexp.Regexp {
	switch {
	case rte.actionPattern != nil && rte.actionPattern.MatchString(r.Path):
		return rte.actionPattern
	case rte.objectPattern != nil && rte.objectPattern.MatchString(r.Path):
		return rte.objectPattern
	case rte.indexPattern.MatchString(r.Path):
		return rte.indexPattern
	}
	return nil
}

func (rte *route) GetParams(r *Request) map[string]string {
	params := make(map[string]string)
	pattern := rte.Match(r)
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

func (rte *route) GetActionName(r *Request) (action string) {
	atIndex := rte.indexPattern.MatchString(r.Path)
	switch {
	case rte.actionPattern != nil && rte.actionPattern.MatchString(r.Path):
		segments := strings.Split(r.Path, "/")
		action = segments[len(segments)-1]
	case atIndex && r.Method == "GET":
		action = "index"
	case atIndex && r.Method == "POST":
		action = "create"
	case !atIndex && r.Method == "GET":
		action = "show"
	case !atIndex && (r.Method == "PUT" || r.Method == "PATCH"):
		action = "update"
	case !atIndex && r.Method == "DELETE":
		action = "destroy"
	}
	return
}

func (rte *route) Respond(r *Request) (status int, body interface{}, action string) {
	action = rte.GetActionName(r)
	if action == "" {
		return 404, "", ""
	}
	r.UrlParams = rte.GetParams(r)
	r.setUser()
	status, body = rte.controller.runFilters(r, action)
	if status != 0 {
		return
	}
	var methodName string
	if extra, ok := rte.controller.extraActions()[action]; ok {
		methodName = extra
	} else {
		methodName = strings.Title(action)
	}
	t := reflect.TypeOf(rte.controller)
	method, _ := t.MethodByName(methodName)
	arguments := []reflect.Value{reflect.ValueOf(rte.controller), reflect.ValueOf(r)}
	statusAndBody := method.Func.Call(arguments)
	status = int(statusAndBody[0].Int())
	body = statusAndBody[1].Interface()
	return
}
