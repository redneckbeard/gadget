package main

import (
	_ "{{.path}}/controllers"
	"github.com/redneckbeard/gadget"
	"github.com/redneckbeard/gadget/templates"
)

type {{.projectName}} struct{}

func (f *{{.projectName}}) OnStart() error {
	gadget.Routes(
		gadget.SetIndex("home"),
	)
	// Delete this line to remove JSON responses
	gadget.Accept("application/json").Via(gadget.JsonBroker)
	// Delete this line to remove template processing for HTML responses
	gadget.Accept("text/html", "*/*").Via(templates.TemplateBroker)
	return nil
}

func main() {
	gadget.SetApp(&{{.projectName}}{})
	gadget.Go()
}
