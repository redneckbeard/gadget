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

type Controller interface {
	Index(r *Request) (int, interface{})
	Show(r *Request) (int, interface{})
	Create(r *Request) (int, interface{})
	Update(r *Request) (int, interface{})
	Destroy(r *Request) (int, interface{})

	IdPattern() string
	ExtraActions() map[string]string
	ExtraActionNames() []string
	setActions([][]string)
	Filter(verbs []string, filter Filter)
	RunFilters(r *Request, action string) (int, interface{})
	Plural() string
}

func NameOf(c Controller) string {
	name := reflect.TypeOf(c).Elem().Name()
	matches := controllerName.FindStringSubmatch(name)
	if matches == nil || len(matches) != 2 {
		panic(`Controller names must adhere to the convention of '<name>Controller'`)
	}
	return strings.ToLower(matches[1])
}

func PluralOf(c Controller) string {
	pluralName := c.Plural()
	if pluralName != "" {
		return pluralName
	}
	return NameOf(c) + "s"
}

func Register(c Controller) {
	c.setActions(arbitraryActions(c))
	controllers[PluralOf(c)] = c
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
