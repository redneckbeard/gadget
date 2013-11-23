package main

import (
	_ "{{.path}}/controllers"
	"{{.path}}/app"
	"github.com/redneckbeard/gadget"
)

func main() {
	gadget.SetApp(app.App)
	gadget.Go()
}
