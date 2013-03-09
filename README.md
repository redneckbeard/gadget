Gadget is a web framework for Go
================================

[![Build Status](https://secure.travis-ci.org/redneckbeard/gadget.png?branch=master)](http://travis-ci.org/redneckbeard/gadget)

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
	_ "inspector/controllers"
)

func main() {
	gadget.Routes(
		gadget.Resource("mission",
			gadget.Resource("character")))

	gadget.Go("8090")
}
```

What's happening here?

* We imported some stuff, notably our controllers with the blank identifier to
  ensure that `init` functions run in that package.
* We configured routes to RESTful controllers and made them nested. We'll now
  have URLs like `/missions/3/characters/7`.
* We configured the server to run on port 8090 with the `Go` function.

Gadget favors convention over configuration (sometimes), and the strings that
are fed to the `gadget.Resource` calls correspond to the names of controllers
that we defined in our `controllers` package. The files in the controller
package all declare a controller type, embed a default controller to make it
simpler to implement the controller interface, and explicitly register that
controller with the framework. Observe:

```Go
package controllers

import (
	"github.com/redneckbeard/gadget"
)

type MissionController struct {
	*gadget.DefaultController
}

func (c *MissionController) Index(r *gadget.Request) (int, interface{}) {
	return 200, []&struct{Mission string}{{"Dr. Claw"},{"M.A.D. Cat"}}
}

func (c *MissionController) Show(r *gadget.Request) (int, interface{}) {
	missionId := r.UrlParams["mission_id"]
	return 200, "Mission #" + missionId + ": this message will self-destruct."
}

func (c *MissionController) ChiefQuimby(r *gadget.Request) (int, interface{}) {
	return 200, "You've done it again, Gadget! Don't know how you do it!"
}

func init() {
        c := &MissionController{gadget.New()}
	gadget.Register(c)
}
```

Controller methods have access to a `Request` object that contains a map of any
parameters pulled from the URL. They return simply an HTTP status code and any
value at all for the body (more on why in a bit). The controller interface
requires `Index`, `Show`, `Create`, `Update`, and `Destroy` methods. Embedding
a pointer to a `DefaultController` means that these are all implemented for
you. However, _this doesn't provide you with anything but 404s_. If you want to
take action in response to a particular verb, override the method.

When we call `controller.Register`, we are then able to `controller.Get` each
of the registered controllers by the lower-cased name of the controller struct
type, minus "controller". But you probably won't ever do so directly, because
the `routing` package does it for you.

The Gadget router will hit controller methods based on the HTTP verbs that you
would expect: 

* `GET /controller` routes to `Index`
* `GET /controller/\d+` routes to `Show`
* `POST /controller` routes `Create`
* `PUT /controller/\d+` routes to `Update`
* `DELETE /controller/\d+` routes to `Destroy`

In addition, any exported method on the controller will be routed to for all
HTTP verbs. `ChiefQuimby` above would be called for any verb when the requested
path was `/mission/chief-quimby`.

Numeric ids are the default, but if you want something else in your URLs, just
override `func IdPattern() string` on your controller.

### Action filters

When developing a web application, you frequently have a short-circuit pattern
common to a number of controller methods -- "404 if the user isn't logged in",
"Redirect if the user isn't authorized", etc. To accommodate code reuse, Gadget
controllers allow you to define filters on certain actions. Setting one up in
the example above might look like this:

```Go
func init() {
	c := &MissionController{gadget.New()}
	c.Filter([]string{"create", "update", "destroy"}, UserIsPenny)
	gadget.Register(c)
}
```

`UserIsPenny` is just a function with the signature `func(r *requests.Request)
(int, interface{})` just like a controller method. If this function returns a
non-zero status code, the controller method that was filtered will never be
called. If the filter returns a status code of zero, Gadget will move on to the
next filter for that action until they are exhausted, and then call the
controller method.

Running Gadget applications
---------------------------

Since your Gadget application is just a Go package, we can build this with `go
install inspector`, and voilà -- we have a single-file web application / HTTP
server waiting for as `$GOPATH/bin/inspector`.

Because there are some files that don't go into the build, and the build is
just an executable, Gadget needs an absolute path that it can assume as the
root that all relative filepaths branch off of. In development, this will often
simply be the current working directory, and that's the default. However, in
production, you might have your binary and your frontend files in completely
different locations. For this reason, we can call the `inspector` executable
with a `-root` flag and point it at whatever path we please.

Gadget assumes that the file root will contain a `static` directory and that
you want it to serve the contents thereof as files. By default, it will do so
at `/static/`. You can, however, change this to accommodate whatever you have
against the word "static"... with the `-static` flag.

The command invoked in an upstart job might then look like:

    /usr/local/bin/inspector -static="/media/" -root=/home/penny/files/

Response processing
-------------------

The interface{} value you that you return from a controller method is by
default piped through `fmt.Sprint`. Strings are predictable, as are numbers;
other types look more like debugging output. However, Gadget has a mechanism
for transforming those values based on `Content-Type` or `Accept` headers. By
defining processor functions and assigning them to MIME types, you can make the
same controller methods speak HTML and JSON.

    processor.Define("application/json", processor.JsonProcessor)
    processor.Define("text/xml", processor.XmlProcessor)

JSON and XML processors are included with Gadget in
`github.com/redneckbeard/processor`. Placing the lines above in your `main`
function will make Gadget serialize the body values returned from your
controller methods when the appropriate headers are found in the request.

Gadget also ships with `processor.TemplateProcessor`, which wraps Go's
excellent html/template package. For any given route, it will load two
templates:

* `$GADGET_ROOT/templates/base.html`, and
* `$GADGET_ROOT/templates/<controller_name>/<action_name>.html`, where
  `controller_name` is the same as the value used in routes, and `action_name`
is the name of the controller method being invoked as a lowercase string.

If either of these templates is not found, the processor will return a 404. If
both are found, the resulting `*Template` will be executed with the interface{}
value returned from the controller method as the context.

This is a thin wrapper around the templating package. There is no abstraction
around the templating language itself, so the controller method template files
must explicitly define a template (or multiple templates) and `base.html` must
reference those templates for the system to actually work.


Wish list
---------

Here are a few features that are clearly necessary but I have yet to implement:

* Methods on gadget.Request for altering response headers, setting cookies, etc.
* Django-style middleware

Thanks for watching
-------------------

![Inspector Gadget](http://www.disneyclips.com/imagesnewb6/imageslwrakr01/inspectorgadget4.gif)
