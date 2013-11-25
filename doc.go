/*
Gadget is a web framework for Go that focuses on building RESTful web services
with content negotiation.

Installation

This overview assumes that you already have a working installation of Go1.1+
and have a defined `$GOPATH`. Just run:

	go get github.com/redneckbeard/gadget

This will install a few subpackages and dependencies. It will not install all
the Gadget subpackages, which provide utilities for form validation, sitemap
generation, and some other stuff.

Project layout

For the most part, you can lay out your Gadget projects however you want. There
is, however, a convention, and you can conform to it most easily by installing
the gadget/gdgt subpackage. Using the "new" command, you can create a
ready-to-compile program with the following directory/file structure:

	.
	├── app
	│   └── conf.go
	├── controllers
	│   └── home.go
	├── main.go
	├── static
	│   ├── css
	│   ├── img
	│   └── js
	└── templates
	    ├── base.html
	    └── home
		└── index.html

The app package is where you actually have your Gadget configuration and a
pointer to the app object, so you will end up importing that package in files
in the controllers package. main.go also imports app, and actually runs the
thing.

Routing

Routes in Gadget are just code. They go in the Configure method of your app in
app/conf.go. 

People love HandlerFuncs. So if you want, you can just route stuff to
HandlerFuncs.

	app.Routes(
		app.HandleFunc("robots.txt", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "templates/robots.txt") }),
	)

If you have lots of routes with a common URL segment, you can factor it out
with Prefixed.

	app.Routes(
		app.Prefixed("users",
			app.HandleFunc("friends", FriendsIndex),
			app.HandleFunc("frenemies", FrenemiesIndex),
		),
	)

In practice, though being able to use HandlerFuncs is very handy, you'll more
commonly route to RESTful controllers with the Resource method.

	app.Routes(
		app.Resource("users",
			app.Resource("friends"),
			app.Resource("frenemies"),
		),
	)

In this example, "users", "friends", and "frenemies" all reference controllers
that we've registered with the app. To mount at route at the root of the site,
there's a special SetIndex method, which is set up for you automatically if you
use the gdgt project generator.

	app.Routes(
		app.SetIndex("home"),
		app.Resource("users",
			app.Resource("friends"),
			app.Resource("frenemies"),
		),
	)

Controllers

The strings that are fed to the `gadget.Resource` calls correspond to the names
of controllers that we defined in our `controllers` package. The files in the
controllers package all declare a gadget.Controller type, embed a pointer to
*gadget.DefaultController to make it simpler to implement the controller
interface, and explicitly register that controller with the framework. 

	package controllers

	import (
		"github.com/redneckbeard/gadget"
		"example/app"
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
		app.Register(&MissionController{})
	}

Controller methods have access to a `Request` object and return simply an HTTP
status code and any value at all for the body (more on why in a bit). The
controller interface requires `Index`, `Show`, `Create`, `Update`, and
`Destroy` methods. Embedding a pointer to a `DefaultController` means that
these are all implemented for you. However, this doesn't provide you with
anything but 404s. If you want to take action in response to a particular verb,
override the method.

The Gadget router will hit controller methods based on the HTTP verbs that you
would expect: 

	`GET /missions` routes to `Index`
	`GET /missions/\d+` routes to `Show`
	`POST /missions` routes `Create`
	`PUT /missions/\d+` routes to `Update`
	`DELETE /missions/\d+` routes to `Destroy`

Numeric ids are the default, but if you want something else in your URLs, just
override `func IdPattern() string` on your controller.

In addition, any exported method on the controller will be routed to for all
HTTP verbs. `ChiefQuimby` above would be called for any verb when the requested
path was `/missions/chief-quimby`.

You make a Controller available to the router by passing it to `app.Register`.
This is best done in the init function. Gadget doesn't pretend to speak
perfect English, so it takes the dumbest possible guess at pluralizing your
controller's name and just tacks an "s" on the end. If inflecting is more
complicated, define a `Plural() string` method on your controller.

Explicit Response objects

In most cases, returning a status code and a response body are all you need to
do to respond to a request. When you do need to set cookies or response
headers, you can wrap the response body value in gadget.NewResponse and set
cookies and headers on the value returned.

Action filters

When developing a web application, you frequently have a short-circuit pattern
common to a number of controller methods -- "404 if the user isn't logged in",
"Redirect if the user isn't authorized", etc. To accommodate code reuse, Gadget
controllers allow you to define filters on certain actions. Setting one up in
the example above might look like this:

	func init() {
		c := &MissionController{gadget.New()}
		c.Filter([]string{"create", "update", "destroy"}, UserIsPenny)
		gadget.Register(c)
	}

`UserIsPenny` is just a function with the signature `func(r *requests.Request)
(int, interface{})` just like a controller method. If this function returns a
non-zero status code, the controller method that was filtered will never be
called. If the filter returns a status code of zero, Gadget will move on to the
next filter for that action until they are exhausted, and then call the
controller method.

Brokers

The interface{} value you that you return from a controller method is by
default piped through `fmt.Sprint`. Strings are predictable, as are numbers;
other types look more like debugging output. However, Gadget has a mechanism
for transforming those values based on `Content-Type` or `Accept` headers. By
defining processor functions and assigning them to MIME types, you can make the
same controller methods speak HTML and JSON.

    app.Accept("application/json").Via(gadget.JsonBroker)

JSON and XML processors are included with Gadget. Placing the line above in
your app's `Configure` method will make Gadget serialize the body values
returned from your controller methods when the appropriate headers are found in
the request.

Subpackage gadget/template implements an HTML Broker that wraps the
html/template package.

Users

Gadget doesn't have an authentication framework, but it does have hooks for
plugging one in. It defines: 1) an interface gadget.User that has a single
method, `Authenticated() bool`; 2) a default anonymous user (Authenticated is
always false); 3) a function type UserIdentifier with the signature
`func(*Request) User`

Once you define your Authenticated User type, you write a UserIdentifier that
will return either your type or gadget.AnonymousUser, and register it with the
framework by passing it to gadget.IdentifyUsersWith.

Debugging

Beyond the debug flag that comes from subpackage `gadget/env` that provides the
"serve" command, you can hook into Request.Debug() on a per-request basis by
setting gadget.SetDebugWith to a function you define with the signature
`func(*Request) bool`.

Asset files

Gadget assumes that the file root will contain a `static` directory and that
you want it to serve the contents thereof as files. By default, it will do so
at `/static/`. You can, however, change this to accommodate whatever you have
against the word "static"... with the `-static` flag.

Running Gadget

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


The command invoked in an upstart job might then look like:

    /usr/local/bin/inspector serve -static="/media/" -root=/home/penny/files/

*/
package gadget
