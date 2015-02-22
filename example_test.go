package gadget

type (
	IndexController  struct{ *DefaultController }
	AuthorController struct{ *DefaultController }
	EntryController  struct{ *DefaultController }
	ExampleApp       struct{ *App }
)

func (ex *ExampleApp) Configure() error {
	ex.Routes(
		ex.SetIndex("index"),
		ex.Prefixed("writing",
			ex.Resource("authors",
				ex.Resource("entries"),
			)))
	return nil
}

var ex *ExampleApp

func (c *IndexController) Plural() string { return "index" }
func (c *EntryController) Plural() string { return "entries" }

func ExampleRoutes() {
	ex = &ExampleApp{&App{}}
	ex.Register(&IndexController{}, &AuthorController{}, &EntryController{})
	ex.Configure()
	ex.printRoutes()
	// Output:
	// ^/(?P<index_id>\d+)$ 						 gadget.IndexController
	// ^writing/authors/(?P<author_id>\d+)$ 				 gadget.AuthorController
	// ^writing/authors/(?P<author_id>\d+)/entries/(?P<entry_id>\d+)$ 	 gadget.EntryController
}
