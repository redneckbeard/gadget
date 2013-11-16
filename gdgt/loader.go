package main

import (
	"os"
	"text/template"
	"path/filepath"
)

func getTemplate(name string) *template.Template {
	path := filepath.Join("$GOPATH", "src/github.com/redneckbeard/gadget/gdgt/templates", name)
	t, err := template.ParseFiles(os.ExpandEnv(path))
	if err != nil {
		panic(err)
	}
	return t
}
