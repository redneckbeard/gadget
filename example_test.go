package gadget_test

import "github.com/redneckbeard/gadget"

type (
	IndexController  struct{ *gadget.DefaultController }
	AuthorController struct{ *gadget.DefaultController }
	EntryController  struct{ *gadget.DefaultController }
)

func (c *IndexController) Plural() string { return "index" }
func (c *EntryController) Plural() string { return "entries" }

func ExampleRoutes() {
	gadget.Register(&IndexController{}, &AuthorController{}, &EntryController{})
	gadget.Routes(
		gadget.SetIndex("index"),
		gadget.Prefixed("writing",
			gadget.Resource("authors",
				gadget.Resource("entries"),
			)))
	gadget.PrintRoutes()
	// Output:
	// ^$ 									 gadget_test.IndexController
	// ^writing/authors(?:/(?P<author_id>\d+))?$ 				 gadget_test.AuthorController
	// ^writing/authors/(?P<author_id>\d+)/entries(?:/(?P<entry_id>\d+))?$ 	 gadget_test.EntryController
}
