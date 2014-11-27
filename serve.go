package gadget

import (
	"github.com/redneckbeard/gadget/env"
	"github.com/redneckbeard/quimby"
	"net/http"
)

var app gdgt

type gdgt interface {
	Configure() error
	Handler() http.HandlerFunc
	Register(...Controller)
	printRoutes()
	GetRoutes() []*route
}

// SetApp registers the top-level app with Gadget.
func SetApp(g gdgt) {
	app = g
}

// Go calls an app's Configure method and runs the command parser.
func Go() {
	if app == nil {
		panic("No call to SetApp found. Ensure that you've imported your app package you are calling SetApp outside of main.")
	}
	if err := app.Configure(); err != nil {
		panic(err)
	}
	env.Handler = app.Handler()
	quimby.Run()
}
