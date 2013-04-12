package gadget

import (
	"github.com/redneckbeard/gadget/env"
	"net/http"
)

var app Gadget

type Gadget interface {
	OnStart() error
}

func SetApp(g Gadget) {
	app = g
}

func Go(port string) {
	if err := env.Configure(); err != nil {
		panic(err)
	}
	if app == nil {
		panic("No call to SetApp found. Ensure that you've imported your app package you are calling SetApp outside of main.")
	}
	if err := app.OnStart(); err != nil {
		panic(err)
	}
	env.ServeStatic()
	http.HandleFunc("/", Handler())
	http.ListenAndServe(":"+port, nil)
}
