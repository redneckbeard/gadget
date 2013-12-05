Gadget is a web framework for Go
================================

[![Build Status](https://travis-ci.org/redneckbeard/gadget.png?branch=master)](https://travis-ci.org/redneckbeard/gadget)

Gadget is a smallish web application framework with a soft spot for content
negotiation. It requires a working installation of Go1.1 and a workspace/
environment variables set up according to [How to Write Go
Code](http://golang.org/doc/code.html).

You can install and all its subpackages at once:

```
go get -v github.com/redneckbeard/gadget/...
```

You can read an [overview of how it works][1], or work through the quick sample
project below.

[1]: http://redneckbeard.github.io/gadget/

Okay. I've got some photos I'd like to share with some
friends, so we'll make a photo app. I'm going to run the gdgt tool at the root
of `$GOPATH/src`.

```
$ gdgt new photos
Created directory photos/controllers
Created directory photos/app
Created directory photos/templates
Created directory photos/static/css
Created directory photos/static/img
Created directory photos/static/js
Created photos/main.go
Created app/conf.go
Created photos/controllers/home.go
```

The skeleton project is empty, but it is a complete program and we can compile it right now. Running `go install photos` will put an executable in `$GOPATH/bin`. I have that directory on my `$PATH`, so now I can invoke `photos`.

```
$ photos
Available commands:
  help
  serve
  list-routes
Type 'photos help <command>' for more information on a specific command.
```

I'm curious what `list-routes` will output.

```
$ photos list-routes
^$ 	 controllers.HomeController
```

So, we get a regular expression representing the root url mapped to
`controllers.HomeController`. Now I want to see what that does, so I'm going to
run `photos serve` and throw a request at it on another tab. I want to build
some super fancy Ember.js client for this site, so I'll politely ask for
JavaScript.

```
$ curl http://127.0.0.1:8090 -H 'Accept: application/json' -v
< HTTP/1.1 404 Not Found
< Content-Type: application/json
< Date: Sat, 23 Nov 2013 05:02:45 GMT
< Content-Length: 2
< 
* Connection #0 to host 127.0.0.1 left intact
""
```

Empty string. Booooo. Well, at least it responded with the correct
`Content-Type`. But we want to send back some sort of content, so we'll open up
`controllers/home.go` and make a few changes. The definition of
`HomeController` looks like this when we open it up:

```
type HomeController struct {
	*gadget.DefaultController
}

func (c *HomeController) Plural() string { return "home" }
```

We're going to add a method to this to serve some content. For now, it's just
going to send back a bunch of strings.

```
func (c *HomeController) Index(r *gadget.Request) (int, interface{}) {
	return 200, []string{
		"pic1.jpg",
		"pic2.jpg",
		"pic3.jpg",
		"pic4.jpg",
	}
}
```

In Gadget, controller methods return a status code and a response body. The
response body's type is `interface{}` -- it can be just about anything you
want. The framework just has to figure out what to do with it.

So we recompile and start the server again, and give that curl another go.

```
$ curl http://127.0.0.1:8090 -H 'Accept: application/json' -v
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Sat, 23 Nov 2013 05:15:18 GMT
< Content-Length: 58
< 
[
  "pic1.jpg",
  "pic2.jpg",
  "pic3.jpg",
  "pic4.jpg"
]
```

That's more like it! But what happens if we don't explicitly request JSON?

```
$ curl http://127.0.0.1:8090
<html>
  <head>
    <title></title>
  </head>
  <body>
	  
[pic1.jpg pic2.jpg pic3.jpg pic4.jpg]

  </body>
</html>
```

That may look a little broken, but it highlights a fundamental principle of
Gadget's controller mechanism: all requests are subject to content negotiation.
This is why the return signature for the response body is `interface{}` -- it
lets us send a single value back to the requester, with the final
representation of that value being determined by a series of content brokers
that match Accept headers against functions for transforming interface values.

The default Gadget configuration sends `*/*` through to a broker that passes
the return value to a Go template. The template we're interested in was created
for us by `gdgt new` as `templates/home/index.html`. It looks like this:

```
{{define "main"}}
{{.}}
{{end}}
```

Not much there. We can change it to loop over the context of the template using
the "range" function, and throw a few tags in there.

```
{{define "main"}}
<ul>
{{range .}}
  <li>{{.}}</li>
{{end}}
</ul>
```

Now our `Accept`less curl looks a little less funky:

```
$ curl http://127.0.0.1:8090
<html>
  <head>
    <title></title>
  </head>
  <body>
  <ul>
    <li>pic1.jpg</li>
    <li>pic2.jpg</li>
    <li>pic3.jpg</li>
    <li>pic4.jpg</li>
  </ul>
  </body>
</html>
```

We're hitting / here, but even if we weren't, there's no extension involved:
Gadget speaks HTTP and only cares about your `Accept` header. We can send
`Accept: text/plain` and we will get an unadorned fmt.Print of the response
body `interface{}` value:

```
$ curl http://127.0.0.1:8090 -H 'Accept: text/plain'
[pic1.jpg pic2.jpg pic3.jpg pic4.jpg]
```

Of course, the data we're getting back could be a bit more interesting. Let's
amend that `Index` method to return a slice of anonymous structs with a few
different fields.

```
func (c *HomeController) Index(r *gadget.Request) (int, interface{}) {
	return 200, []struct{
		Filename, Title, Description string
	}{
		{"pic1.jpg", "Pic One", "The first picture in our albom."},
		{"pic2.jpg", "Pic Two", "The second picture in our albom."},
		{"pic3.jpg", "Pic Three", "The third picture in our albom."},
		{"pic4.jpg", "Pic Four", "The fourth picture in our albom."},
	}
}
```

We'll also make a tweak to the inside of that template loop.

```
<li>
  <a href="{{.Filename}}" title="{{.Title}}">{{.Description}}</a>
</li>
```

Now our responses for the three content types we've tried sending look like this:

```
$ curl http://127.0.0.1:8090 -H 'Accept: text/plain'
[{pic1.jpg Pic One The first picture in our albom.} {pic2.jpg Pic Two The second picture in our albom.} {pic3.jpg Pic Three The third picture in our albom.} {pic4.jpg Pic Four The fourth picture in our albom.}]

jlukens-mbp:src jlukens$ curl http://127.0.0.1:8090 -H 'Accept: text/html'
<html>
  <head>
    <title></title>
  </head>
  <body>
  <ul>
    <li>
      <a href="pic1.jpg" title="Pic One">The first picture in our albom.</a>
    </li>
    <li>
      <a href="pic2.jpg" title="Pic Two">The second picture in our albom.</a>
    </li>
    <li>
      <a href="pic3.jpg" title="Pic Three">The third picture in our albom.</a>
    </li>
    <li>
      <a href="pic4.jpg" title="Pic Four">The fourth picture in our albom.</a>
    </li>
  </ul>
  </body>
</html>

$ curl http://127.0.0.1:8090 -H 'Accept: application/json'
[
  {
    "Filename": "pic1.jpg",
    "Title": "Pic One",
    "Description": "The first picture in our albom."
  },
  {
    "Filename": "pic2.jpg",
    "Title": "Pic Two",
    "Description": "The second picture in our albom."
  },
  {
    "Filename": "pic3.jpg",
    "Title": "Pic Three",
    "Description": "The third picture in our albom."
  },
  {
    "Filename": "pic4.jpg",
    "Title": "Pic Four",
    "Description": "The fourth picture in our albom."
  }
]
```

This is an extremely brief introduction, but if you think it's neat and would
like to know more, there are plenty of docs [over on
godoc.org](http://godoc.org/github.com/redneckbeard/gadget).
