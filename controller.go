package gadget

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

func init() {
	controllerName = regexp.MustCompile(`(\w+)Controller`)
	pascalCase = regexp.MustCompile(`[A-Z]+[a-z\d]+`)
}

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
	pascalCase     *regexp.Regexp
	defaultActions = []string{INDEX, SHOW, CREATE, UPDATE, DESTROY}
)

// Controller is the interface that defines how Gadget applications respond to
// HTTP requests at a given URL with a particular HTTP verb. Controllers have five
// primary methods for handling requests. For a Controller mounted at /foo,
// requests are routed to methods as follows:
// 
// 	GET 	/foo 			Index
// 	POST 	/foo			Create
// 	GET 	/foo/<idPattern> 	Show
// 	PUT 	/foo/<idPattern>	Update
// 	DELETE 	/foo/<idPattern> 	Destroy
// 
// Each of these methods takes a *gadget.Request as its only argument and returns
// an HTTP status code as an int and an interface{} value as the body. The
// interface{} value is then cast to a string, serialized, used as a template
// context, etc. according to the application's broker configuration. 
// 
// Any other exported method with a signature of (*gadget.Request) (int,
// interface{}) will also be routable. For example, if the controller mounted
// above at /foo defines a method AllTheThings(r *gadget.Request) (int,
// interface{}), Gadget will route any request for /foo/all-the-things, regardless
// of verb, to that method.
// 
// Controller also requires two methods that enable users to customize routing
// options to this controller, IdPattern and Plural.  The final exported method of
// the Controller interface is Filter, which allows for abstracting common
// patterns from multiple Controller methods. All three of these methods are
// documented in the fallback implementations provided by DefaultController.
//
// Applications must inform Gadget of the existence of Controller types using
// the Register function.
type Controller interface {
	Index(r *Request) (int, interface{})
	Show(r *Request) (int, interface{})
	Create(r *Request) (int, interface{})
	Update(r *Request) (int, interface{})
	Destroy(r *Request) (int, interface{})

	Filter(verbs []string, filter Filter)
	IdPattern() string
	Plural() string

	extraActionNames() []string
	extraActions() map[string]string
	runFilters(r *Request, action string) (int, interface{})
	setActions([][]string)
}

func nameOf(c Controller) string {
	name := reflect.TypeOf(c).Elem().Name()
	matches := controllerName.FindStringSubmatch(name)
	if matches == nil || len(matches) != 2 {
		panic(`Controller names must adhere to the convention of '<name>Controller'`)
	}
	return strings.ToLower(matches[1])
}

func pluralOf(c Controller) string {
	pluralName := c.Plural()
	if pluralName != "" {
		return pluralName
	}
	return nameOf(c) + "s"
}

// Register notifies Gadget that you want to use a type as a Controller. It takes
// a variable number of arguments, all of which must be pointers to struct types
// that embed DefaultController. It will panic if the name of the Controller is
// not in the form <resource name>Controller or if a struct value is passed
// instead of a pointer.
// 
// Register makes naive assumptions about the name of your Controller; firstly,
// that it is English; secondly that it is in singular number; and thirdly, that
// it can be pluralized simply by appending the letter "s". Thus,
// "WidgetController" will be available to the Resource function for routes as
// "widgets", but "EntryController" as "entrys", unless you define the Plural
// method on the Controller to return something more correct.
func Register(clist ...Controller) {
	for _, c := range clist {
		v := reflect.ValueOf(c).Elem()
		defaultCtlr := v.FieldByName("DefaultController")
		defaultCtlr.Set(reflect.ValueOf(newController()))
		c.setActions(arbitraryActions(c))
		controllers[pluralOf(c)] = c
	}
}

func clear() {
	controllers = make(map[string]Controller)
}

func getController(name string) (Controller, error) {
	controller, ok := controllers[name]
	if !ok {
		return nil, errors.New(fmt.Sprintf("No controller with label '%s' found", name))
	}
	return controller, nil
}

func hyphenate(pascal string) string {
	matches := []string{}
	for _, match := range pascalCase.FindAllString(pascal, -1) {
		matches = append(matches, strings.ToLower(match))
	}
	return strings.Join(matches, "-")
}

func arbitraryActions(ctlr Controller) (actions [][]string) {
	t := reflect.TypeOf(ctlr)
	v := reflect.ValueOf(ctlr)
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		if method.PkgPath == "" && isAction(v.Method(i)) {
			isDefault := false
			name := hyphenate(method.Name)
			for _, da := range defaultActions {
				if name == da {
					isDefault = true
				}
			}
			if !isDefault {
				actions = append(actions, []string{name, method.Name})
			}
		}
	}
	return
}

func isAction(method reflect.Value) bool {
	indexMethod := reflect.ValueOf(&DefaultController{}).MethodByName("Index")
	return indexMethod.Type() == method.Type()
}
