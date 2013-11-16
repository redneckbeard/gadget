package controllers

import (
	"github.com/redneckbeard/gadget"
)

func init() {
	gadget.Register(&{{.name}}Controller{})
}

type {{.name}}Controller struct {
	*gadget.DefaultController
}
