package controllers

import (
	"github.com/redneckbeard/gadget"
	"{{.project}}/app"
)

func init() {
	app.Register(&{{.name}}Controller{})
}

type {{.name}}Controller struct {
	*gadget.DefaultController
}

{{if matchString .name "Home"}}
func (c *HomeController) Plural() string { return "home" }
{{end}}
