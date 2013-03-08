package controller

import (
	"fmt"
	"github.com/redneckbeard/gadget/requests"
	"reflect"
	"strings"
)

const (
	INDEX   = "index"
	SHOW    = "show"
	CREATE  = "create"
	UPDATE  = "update"
	DESTROY = "destroy"
)

var defaultActions = []string{INDEX, SHOW, CREATE, UPDATE, DESTROY}

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
	Index(r *requests.Request) (int, interface{})
	Show(r *requests.Request) (int, interface{})
	Create(r *requests.Request) (int, interface{})
	Update(r *requests.Request) (int, interface{})
	Destroy(r *requests.Request) (int, interface{})

	IdPattern() string
	ExtraActions() []string
	setActions([]string)
	Filter(verbs []string, filter Filter)
	RunFilters(r *requests.Request, action string) (int, interface{})
}

func New() *DefaultController {
	return &DefaultController{}
}

type DefaultController struct {
	filters      map[string][]Filter
	extraActions []string
}

// Default handlers for anything just 404. Users are responsible for overriding any embedded methods they want to respond otherwise.
func (c *DefaultController) Index(r *requests.Request) (int, interface{})   { return 404, "" }
func (c *DefaultController) Show(r *requests.Request) (int, interface{})    { return 404, "" }
func (c *DefaultController) Create(r *requests.Request) (int, interface{})  { return 404, "" }
func (c *DefaultController) Update(r *requests.Request) (int, interface{})  { return 404, "" }
func (c *DefaultController) Destroy(r *requests.Request) (int, interface{}) { return 404, "" }

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

func (c *DefaultController) RunFilters(r *requests.Request, action string) (status int, body interface{}) {
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
