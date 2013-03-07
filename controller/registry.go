package controller

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var (
	controllers    = make(map[string]Controller)
	controllerName *regexp.Regexp
)

func init() {
	controllerName = regexp.MustCompile(`(\w+)Controller`)
}

func NameOf(c Controller) string {
	name := reflect.TypeOf(c).Elem().Name()
	matches := controllerName.FindStringSubmatch(name)
	if matches == nil || len(matches) != 2 {
		panic("Controller names must adhere to the convention of '<name>Controller'")
	}
	return strings.ToLower(matches[1])
}

func Register(c Controller) {
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
