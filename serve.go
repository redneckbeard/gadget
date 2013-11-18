package gadget

import (
	"github.com/redneckbeard/gadget/env"
	"github.com/redneckbeard/quimby"
	"net/http"
)

var app Gadget

type Gadget interface {
	Configure() error
	Handler() http.HandlerFunc
	PrintRoutes()
	Register(...Controller)
}

func SetApp(g Gadget) {
	app = g
}

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
