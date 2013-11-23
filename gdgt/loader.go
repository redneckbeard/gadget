package main

import (
	"os"
	"text/template"
	"path/filepath"
	"io"
	"fmt"
)

func templatePath(name string) string {
	return os.ExpandEnv(filepath.Join("$GOPATH", "src/github.com/redneckbeard/gadget/gdgt/templates", name))
}

func getTemplate(name string) *template.Template {
	path := templatePath(name)
	t, err := template.ParseFiles(path)
	if err != nil {
		panic(err)
	}
	return t
}

func copyTemplate(src, dst string) {
	templateFile, err := os.Open(templatePath(src))
	if err != nil {
		fmt.Println("Failed to open template file " + templatePath(src))
		fmt.Println(err)
	}
	projectFile, err := os.Create(dst)
	io.Copy(projectFile, templateFile)
}
