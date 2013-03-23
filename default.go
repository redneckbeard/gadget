package gadget

import "fmt"

func New() *DefaultController {
	controller := &DefaultController{}
	controller.filters = make(map[string][]Filter)
	controller.extraActions = make(map[string]string)
	return controller
}

type DefaultController struct {
	filters      map[string][]Filter
	extraActions map[string]string
}

type Filter func(*Request) (int, interface{})

// Default handlers for anything just 404. Users are responsible for overriding any embedded methods they want to respond otherwise.
func (c *DefaultController) Index(r *Request) (int, interface{})   { return 404, "" }
func (c *DefaultController) Show(r *Request) (int, interface{})    { return 404, "" }
func (c *DefaultController) Create(r *Request) (int, interface{})  { return 404, "" }
func (c *DefaultController) Update(r *Request) (int, interface{})  { return 404, "" }
func (c *DefaultController) Destroy(r *Request) (int, interface{}) { return 404, "" }

func (c *DefaultController) IdPattern() string { return `\d+` }
func (c *DefaultController) Plural() string    { return "" }

func (c *DefaultController) Filter(verbs []string, filter Filter) {
	for _, verb := range verbs {
		if filters, ok := c.filters[verb]; !ok {
			panic(fmt.Sprintf("Unable to add filter for '%s' -- no such action", verb))
		} else {
			c.filters[verb] = append(filters, filter)
		}
	}
}

func (c *DefaultController) RunFilters(r *Request, action string) (status int, body interface{}) {
	for _, f := range c.filters[action] {
		status, body = f(r)
		if status == 0 {
			return
		}
	}
	return
}

func (c *DefaultController) ExtraActions() map[string]string {
	return c.extraActions
}

func (c *DefaultController) ExtraActionNames() (names []string) {
	for k, _ := range c.extraActions {
		names = append(names, k)
	}
	return
}

func (c *DefaultController) setActions(actions [][]string) {
	for _, action := range actions {
		c.extraActions[action[0]] = action[1]
		c.filters[action[0]] = []Filter{}
	}
	for _, action := range defaultActions {
		c.filters[action] = []Filter{}
	}
}
