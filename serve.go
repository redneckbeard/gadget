package gadget

import (
	"github.com/redneckbeard/gadget/env"
	"github.com/redneckbeard/quimby"
)

var app Gadget

type Gadget interface {
	OnStart() error
}

func SetApp(g Gadget) {
	app = g
}

func Go() {
	if app == nil {
		panic("No call to SetApp found. Ensure that you've imported your app package you are calling SetApp outside of main.")
	}
	if err := app.OnStart(); err != nil {
		panic(err)
	}
	env.Handler = Handler()
	quimby.Run()
}
