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

{{if eq .name "Home"}}
func (c *HomeController) Plural() string { return "home" }
{{end}}
