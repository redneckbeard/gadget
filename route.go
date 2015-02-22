package gadget

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

type route struct {
	segment                                              string
	segments                                             []*segment
	indexPattern, objectPattern, actionPattern, hostname *regexp.Regexp
	handler                                              http.HandlerFunc
	controller                                           Controller
	subroutes                                            []*route
	isRoot                                               bool
}

func (rte *route) String() string {
	if rte.objectPattern != nil {
		return rte.objectPattern.String()
	}
	return rte.indexPattern.String()
}

type segmentList []*segment

func (sl segmentList) patternComponents() (*segment, []string) {
	allButLast := sl[:len(sl)-1]
	segments := []string{}
	for _, s := range allButLast {
		if s.isPrefix {
			segments = append(segments, s.name)
		} else {
			segments = append(segments, s.name, s.objectSuffix())
		}
	}
	return sl[len(sl)-1], segments
}

func (sl segmentList) indexPattern() *regexp.Regexp {
	final, segments := sl.patternComponents()
	return finalPattern(segments, final.name)
}

func (sl segmentList) objectPattern() *regexp.Regexp {
	final, segments := sl.patternComponents()
	return finalPattern(segments, final.name, final.objectSuffix())
}

func (sl segmentList) actionPattern() *regexp.Regexp {
	final, segments := sl.patternComponents()
	return finalPattern(segments, final.name, final.actionSuffix())
}

func finalPattern(segments []string, suffixes ...string) *regexp.Regexp {
	segments = append(segments, suffixes...)
	return regexp.MustCompile(fmt.Sprintf("^%s$", strings.Join(segments, "/")))
}

type segment struct {
	name, paramName, idPattern string
	isPrefix                   bool
	actions                    []string
}

func (s *segment) objectSuffix() string {
	return fmt.Sprintf("(?P<%s_id>%s)", s.paramName, s.idPattern)
}

func (s *segment) actionSuffix() string {
	return fmt.Sprintf("(?:%s)", strings.Join(s.actions, "|"))
}

func (rte *route) buildPatterns(prefix string, segments ...*segment) {
	if rte.controller != nil {
		rte.segments = append(segments, &segment{
			name:      rte.segment,
			paramName: strings.Replace(nameFromController(rte.controller), "-", "_", -1),
			idPattern: rte.controller.IdPattern(),
			actions:   rte.controller.extraActionNames(),
		})
	} else {
		rte.segments = append(segments, &segment{
			name:     prefix,
			isPrefix: true,
		})
	}
	for _, r := range rte.subroutes {
		r.buildPatterns(prefix, rte.segments...)
	}
	patterns := segmentList(rte.segments)
	if rte.controller != nil {
		rte.indexPattern = patterns.indexPattern()
		rte.objectPattern = patterns.objectPattern()
		if len(rte.controller.extraActionNames()) > 0 {
			rte.actionPattern = patterns.actionPattern()
		}
	}
}

func (rte *route) flatten() []*route {
	var flattened []*route
	if rte.controller != nil || rte.handler != nil {
		flattened = append(flattened, rte)
	}
	for _, r := range rte.subroutes {
		r.hostname = rte.hostname
		flattened = append(flattened, r.flatten()...)
	}
	return flattened
}

func (a *App) newRoute(segment string, handler http.HandlerFunc) *route {
	if segment == "" {
		return &route{segment: segment}
	}
	controller, err := a.getController(segment)
	if err != nil && handler == nil {
		panic(err)
	}
	rte := &route{segment: segment, controller: controller, handler: handler}
	return rte
}

func (rte *route) Match(r *Request) *regexp.Regexp {
	host := strings.Split(r.Host, ":")[0]
	if rte.hostname != nil && !rte.hostname.MatchString(host) {
		return nil
	}
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
