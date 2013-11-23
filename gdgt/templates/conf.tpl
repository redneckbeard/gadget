package app

import (
	"github.com/redneckbeard/gadget"
	"github.com/redneckbeard/gadget/templates"
)

var App *gdgt

func init() {
	App = &gdgt{ &gadget.App{} }
}

type gdgt struct { *gadget.App }

// Configure your app here.
func (app *gdgt) Configure() error {
	app.Routes(
		app.SetIndex("home"),
	)
	app.Accept("application/json").Via(gadget.JsonBroker)
	app.Accept("text/html", "*/*").Via(templates.TemplateBroker)
	return nil
}

// Just a little sugar to make common method calls more succinct

func Register(controllers ...gadget.Controller) {
	App.Register(controllers...)
}
