package gadget

import "fmt"

func newController() *DefaultController {
	controller := &DefaultController{}
	controller.filters = make(map[string][]Filter)
	controller.extraActionMap = make(map[string]string)
	return controller
}

// DefaultController satisfies the Controller interface, and because Controller
// includes unexported methods, must be embedded in any other type that implements
// Controller. The fallback implementations provided by DefaultController and the
// interactions of other Gadget machinery therewith can be summarized as follows:
// 
// 	* Index, Show, Create, Update, and Destroy methods all return a 404 with 
// 	  an empty string for the body
// 	* The return value of IdPattern for use in routes is `\d+`
// 	* The return value of Plural is "", which Register takes to mean "just 
// 	  add an 's'"
type DefaultController struct {
	filters        map[string][]Filter
	extraActionMap map[string]string
}

// Filter is simply a function with the same signature as a controller method
// minus the receiver. They are used in calls to Controller.Filter.
type Filter func(*Request) (int, interface{})

// The default return value of Index is (404, "").
func (c *DefaultController) Index(r *Request) (int, interface{}) { return 404, "" }

// The default return value of Show is (404, "").
func (c *DefaultController) Show(r *Request) (int, interface{}) { return 404, "" }

// The default return value of Create is (404, "").
func (c *DefaultController) Create(r *Request) (int, interface{}) { return 404, "" }

// The default return value of Update is (404, "").
func (c *DefaultController) Update(r *Request) (int, interface{}) { return 404, "" }

// The default return value of Destroy is (404, "").
func (c *DefaultController) Destroy(r *Request) (int, interface{}) { return 404, "" }

// IdPattern returns a string that will be used in a regular expression to
// match a unique identifier for a resource in a URL. The matched value is then
// added to gadget.Request.UrlParams.
func (c *DefaultController) IdPattern() string { return `\d+` }

// Plural returns a string to be used as the plural form of the first word in
// the name of the Controller type. Note that this value needs to be the entire
// plural form of the word, not an ending.
func (c *DefaultController) Plural() string { return "" }

// Filter applies a Filter func to the Controller methods named in the string
// slice verbs. The strings in the slice should be the lowercased name of the
// default method, or for additional methods, the hyphenated string that
// appears in routes.
//
// If any Filter returns non-zero number for the HTTP status, the body of the
// Controller method will never be executed and the response cycle will begin
// with the Filter.
//
// 	c := &PostController{}
// 	c.Filter([]string{"create", "update", "destroy"}, func(r *gadget.Request) (int, interface{}) {
// 		if !r.User.Authenticated() {
// 			return 403, "Verboten"
// 		}
// 	}
// 	gadget.Register(c)
func (c *DefaultController) Filter(verbs []string, filter Filter) {
	for _, verb := range verbs {
		if filters, ok := c.filters[verb]; !ok {
			panic(fmt.Sprintf("Unable to add filter for '%s' -- no such action", verb))
		} else {
			c.filters[verb] = append(filters, filter)
		}
	}
}

func (c *DefaultController) runFilters(r *Request, action string) (status int, body interface{}) {
	for _, f := range c.filters[action] {
		status, body = f(r)
		if status == 0 {
			return
		}
	}
	return
}

func (c *DefaultController) extraActions() map[string]string {
	return c.extraActionMap
}

func (c *DefaultController) extraActionNames() (names []string) {
	for k, _ := range c.extraActionMap {
		names = append(names, k)
	}
	return
}

func (c *DefaultController) setActions(actions [][]string) {
	for _, action := range actions {
		c.extraActionMap[action[0]] = action[1]
		c.filters[action[0]] = []Filter{}
	}
	for _, action := range defaultActions {
		c.filters[action] = []Filter{}
	}
}
