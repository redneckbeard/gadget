Gadget is a web framework for Go
================================

tl;dr I'm a terrible maintainer of open source software, and I would advise
against using this framework.

Installation
------------

This overview assumes that you already have a working installation of Go1 and
have a defined `$GOPATH`.

Gadget includes a number of subpackages. You might not need all of them once
you know how this all works, but for now, `go get` them all at once with the
handy `...` syntax:

    go get -v github.com/redneckbeard/gadget/...

`-v` will display the names of the packages as they are fetched and compiled,
which makes me feel more confident that the command actually did what I wanted.

Project setup
-------------

There's very little imposed file/directory structure in a Gadget application. A
suggested directory layout might look like:

    ./inspector
    ├── main.go
    ├── controllers
    │   ├── missions.go
    │   └── characters.go
    └── static
        ├── css
        ├── img
        └── js

The only fixed piece of this layout is that the directory containing static
assets must be named "static". I cannot imagine a good case for making this
configurable.

`main.go` could look something like this:

```Go
package main

import (
	"github.com/redneckbeard/gadget"
	"github.com/redneckbeard/gadget/routing"
	_ "inspector/controllers"
)

func main() {
	routing.Routes(
		routing.Resource("mission",
			routing.Resource("character")))

	gadget.Go("8090")
}
```

What's happening here?

* We imported some stuff, notably our controllers with the blank identifier to ensure that `init` functions run in that package.
* We configured routes to RESTful controllers and made them nested. We'll now have URLs like `/missions/3/characters/7`.
* We configured the server to run on port 8090 with the `Go` function.

Gadget favors convention over configuration (sometimes), and the strings that are fed to the `routing.Resource` calls correspond to the names of controllers that we defined in our `controllers` package. The files in the controller package all declare a controller type, embed a default controller to make it simpler to implement the controller interface, and explicitly register that controller with the framework. Observe:

```Go
package controllers

import (
	"github.com/redneckbeard/gadget/controller"
	"github.com/redneckbeard/gadget/requests"
)

type MissionController struct {
	*controller.DefaultController
}

func (c *MissionController) Index(r *requests.Request) (int, interface{}) {
	return 200, "I'm a list of foos"
}

func (c *MissionController) Show(r *requests.Request) (int, interface{}) {
	missionId := r.UrlParams["mission_id"]
	return 200, "Mission #" + missionId + ": this message will self-destruct."
}

func init() {
	controller.Register(&FooController{controller.New()})
}
```

Controller methods have access to a `Request` object that contains a map of any parameters pulled from the URL. They return simply an HTTP status code and any value at all for the body (more on why in a bit). The controller interface requires `Index`, `Show`, `Create`, `Update`, and `Destroy` methods. Embedding a pointer to a `DefaultController` means that these are all implemented for you. However, _this doesn't provide you with anything but 404s_. If you want to take action in response to a particular verb, override the method.

When we call `controller.Register`, we are then able to `controller.Get` each
of the registered controllers by the lower-cased name of the controller struct
type, minus "controller". But you probably won't ever do so directly, because
the `routing` package does it for you.

Running Gadget applications
---------------------------

Since your Gadget application is just a Go package, we can build this with `go install inspector`, and voilà -- we have a single-file web application / HTTP server waiting for as `$GOPATH/bin/inspector`.

Because there are some files that don't go into the build, and the build is just an executable, Gadget needs an absolute path that it can assume as the root that all relative filepaths branch off of. In development, this will often simply be the current working directory, and that's the default. However, in production, you might have your binary and your frontend files in completely different locations. For this reason, we can call the `inspector` executable with a `-root` flag and point it at whatever path we please.

Gadget assumes that the file root will contain a `static` directory and that you want it to serve the contents thereof as files. By default, it will do so at `/static/`. You can, however, change this to accommodate whatever you have against the word "static"... with the `-static` flag.

The command invoked in an upstart job might then look like:

    /usr/local/bin/inspector -static="/media/" -root=/home/penny/files/

Response processing
-------------------

The interface{} value you that you return from a controller method is by
default piped through `fmt.Sprint`. Strings are predictable, as are numbers;
other types look more like debugging output. However, Gadget has a mechanism
for transforming those values based on `Content-Type` or `Accept` headers. By defining processor functions and assigning them to MIME types, you can make the same controller methods speak HTML and JSON.

    processor.Define("application/json", processor.JsonProcessor)
    processor.Define("text/xml", processor.XmlProcessor)

JSON and XML processors are included with Gadget in `github.com/redneckbeard/processor`. Placing the lines above in your `main` function will make Gadget serialize the body values returned from your controller methods when the appropriate headers are found in the request.

Wish list
---------

Here are a few features that are clearly necessary but I have yet to implement:

* Ability to override numeric ids in URLs with whatever regexp you want
* A Response type for altering headers, setting cookies, etc.
* Django-style middleware
* A Processor function to handle rendering to templates based on controller/action names à la Rails
* Arbitrary exporteds methods on controllers working as actions

Thanks for watching
-------------------

![Inspector Gadget](http://www.disneyclips.com/imagesnewb6/imageslwrakr01/inspectorgadget4.gif)
