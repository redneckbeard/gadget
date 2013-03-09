package gadget

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

const (
	INDEX   = "index"
	SHOW    = "show"
	CREATE  = "create"
	UPDATE  = "update"
	DESTROY = "destroy"
)

var (
	controllers    = make(map[string]Controller)
	controllerName *regexp.Regexp
	defaultActions = []string{INDEX, SHOW, CREATE, UPDATE, DESTROY}
)

func init() {
	controllerName = regexp.MustCompile(`(\w+)Controller`)
}

func NameOf(c Controller) string {
	name := reflect.TypeOf(c).Elem().Name()
	matches := controllerName.FindStringSubmatch(name)
	if matches == nil || len(matches) != 2 {
		panic(`Controller names must adhere to the convention of '<name>Controller'`)
	}
	return strings.ToLower(matches[1])
}

func Register(c Controller) {
	c.setActions(arbitraryActions(c))
	controllers[NameOf(c)] = c
}

func Clear() {
	controllers = make(map[string]Controller)
}

func Get(name string) (Controller, error) {
	controller, ok := controllers[name]
	if !ok {
		return nil, errors.New(fmt.Sprintf("No controller with label '%s' found", name))
	}
	return controller, nil
}

func arbitraryActions(ctlr Controller) (actions []string) {
	t := reflect.TypeOf(ctlr)
	v := reflect.ValueOf(ctlr)
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		if method.PkgPath == "" && isAction(v.Method(i)) {
			isDefault := false
			name := strings.ToLower(method.Name)
			for _, da := range defaultActions {
				if method.Name == da {
					isDefault = true
				}
			}
			if !isDefault {
				actions = append(actions, name)
			}
		}
	}
	return
}

func isAction(method reflect.Value) bool {
	indexMethod := reflect.ValueOf(&DefaultController{}).MethodByName("Index")
	return indexMethod.Type() == method.Type()
}

type Controller interface {
	Index(r *Request) (int, interface{})
	Show(r *Request) (int, interface{})
	Create(r *Request) (int, interface{})
	Update(r *Request) (int, interface{})
	Destroy(r *Request) (int, interface{})

	IdPattern() string
	ExtraActions() []string
	setActions([]string)
	Filter(verbs []string, filter Filter)
	RunFilters(r *Request, action string) (int, interface{})
}

func New() *DefaultController {
	return &DefaultController{}
}

type DefaultController struct {
	filters      map[string][]Filter
	extraActions []string
}

// Default handlers for anything just 404. Users are responsible for overriding any embedded methods they want to respond otherwise.
func (c *DefaultController) Index(r *Request) (int, interface{})   { return 404, "" }
func (c *DefaultController) Show(r *Request) (int, interface{})    { return 404, "" }
func (c *DefaultController) Create(r *Request) (int, interface{})  { return 404, "" }
func (c *DefaultController) Update(r *Request) (int, interface{})  { return 404, "" }
func (c *DefaultController) Destroy(r *Request) (int, interface{}) { return 404, "" }

func (c *DefaultController) IdPattern() string {
	return `\d+`
}

func (c *DefaultController) Filter(verbs []string, filter Filter) {
	for _, verb := range verbs {
		if filters, ok := c.filters[verb]; !ok {
			panic(fmt.Sprintf("Unable to add filter for '%s' -- no such action", verb))
		} else {
			c.filters[verb] = append(filters, filter)
		}
	}
}

type Filter func(*Request) (int, interface{})

func (c *DefaultController) RunFilters(r *Request, action string) (status int, body interface{}) {
	for _, f := range c.filters[action] {
		status, body = f(r)
		if status == 0 {
			return
		}
	}
	return
}

func (c *DefaultController) ExtraActions() []string {
	return c.extraActions
}

func (c *DefaultController) setActions(actions []string) {
	c.extraActions = actions
	actions = append(actions, defaultActions...)
	filters := make(map[string][]Filter)
	for _, action := range actions {
		filters[action] = []Filter{}
	}
	c.filters = filters
}
