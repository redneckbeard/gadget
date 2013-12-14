package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"
)

func templatePath(name string) string {
	return os.ExpandEnv(filepath.Join("$GOPATH", "src/github.com/redneckbeard/gadget/gdgt/templates", name))
}

func getTemplate(name string) *template.Template {
	path := templatePath(name)
	t, err := template.New(name).Funcs(template.FuncMap{
		"matchString": func(a, b string) bool { return a == b },
	}).ParseFiles(path)
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
	fmt.Println("Created " + dst)
}
