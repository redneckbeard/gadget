HEY CHUCK!
==========

I put this thing I'm working on up on the Githubz. Let's see if you can get it to work.

    go get -v github.com/redneckbeard/gadget/...

Let's say you have a `goapp` folder on your `$GOPATH`. Inside that folder, you
might have a `controllers/foo.go` that looks like this:

	```Go
	package controllers

	import (
		"gadget/controller"
		"gadget/requests"
	)

	type FooController struct {
		// Embed this to implement the Controller interface and get some default methods
		*controller.DefaultController
	}

	func (c *FooController) Index(r *requests.Request) (int, interface{}) {
		// Controller methods return a status code an whatever you want the body to be
		return 200, "I'm a list of foos"
	}

	func (c *FooController) Show(r *requests.Request) (int, interface{}) {
		fooId := r.UrlParams["foo_id"]
		return 200, "I'm foo #" + fooId
	}

	func init() {
		// We have to embed a pointer to a DefaultController to properly initialize a controller for the registry
		controller.Register(&FooController{controller.New()})
	}
	```



We also have some other controllers in `controllers/bar.go` and
`controllers/baz.go`. They happen to be named `BarController` and
`BazController`.

When we call `controller.Register`, we are then able to `controller.Get` each
of the registered controllers by the lower-cased name of the controller struct
type, minus "controller". But you probably won't ever do so directly, because
the `routing` package does it for you.

Our `goapp/main.go` might be:

	```Go
	package main

	import (
		"gadget"
		"gadget/routing"
		_ "goapp/controllers"
	)

	func main() {
		routing.Routes(
			routing.SetIndex("bar"), 
			routing.Resource("foo",
				routing.Resource("baz")),
			routing.Prefixed("admin", 
				routing.Resource("foo"),
				routing.Resource("bar")))
		gadget.Go("8090")
	}
	```

The routing stuff has tests with comments that should explain the URLs this
generates if it's fuzzy.

`gadget.Go` is at this point just a convenience method for wrapping the
standard `net/http` calls to run a Handler on a designated port. If you don't
to use it, you don't have to -- `routing.Handler` returns a valid Handler
function so you can use the gadget controllers/routes alongside other code
using `net/http`.

Additional features planned:

* The `requests` package is going to have a `Response` type that you can return
as the body for playing with headers, cookies, etc.
* Response processor definitions (this is the reason that you return an
`interface{}` value instead of just a string from controller methods) --
there will be a mechanism for per-mimetype plugins for transforming the
`interface{}` body returned into a string. I intend to provide defaults for
`text/html` that will try rendering a template if the value is a struct type or
a map and for `application/json` that will just pipe it through `json.Marshal`.

![Inspector Gadget](http://www.disneyclips.com/imagesnewb6/imageslwrakr01/inspectorgadget4.gif)

Hope you're loving Hacker School!


