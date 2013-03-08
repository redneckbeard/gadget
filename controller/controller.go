package controller

import (
	"fmt"
	"github.com/redneckbeard/gadget/requests"
	"reflect"
)

const (
	INDEX   = "index"
	SHOW    = "show"
	CREATE  = "create"
	UPDATE  = "update"
	DESTROY = "destroy"
)

var VALID_METHOD_NAMES = []string{INDEX, SHOW, CREATE, UPDATE, DESTROY}

type Controller interface {
	Index(r *requests.Request) (int, interface{})
	Show(r *requests.Request) (int, interface{})
	Create(r *requests.Request) (int, interface{})
	Update(r *requests.Request) (int, interface{})
	Destroy(r *requests.Request) (int, interface{})

	IdPattern() string
	Filter(verbs []string, filter Filter)
	RunFilters(r *requests.Request, action string) (int, interface{})
}

func New() *DefaultController {
	filters := make(map[string][]Filter)
	for _, m := range VALID_METHOD_NAMES {
		filters[m] = []Filter{}
	}
	return &DefaultController{filters}
}

type DefaultController struct {
	filters map[string][]Filter
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
		filters := c.filters[verb]
		if filters == nil {
			panic(fmt.Sprintf("Cannot register beforeFilter for verb '%s'", verb))
		}
		c.filters[verb] = append(filters, filter)
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

func IsAction(method reflect.Value) bool {
	indexMethod := reflect.ValueOf(&DefaultController{}).MethodByName("Index")
	return indexMethod.Type() == method.Type()
}
