package main

import (
	"fmt"
	"github.com/redneckbeard/gadget/strutil"
	"github.com/redneckbeard/quimby"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	quimby.Add(&Controller{})
}

// Controller provides a command that generates controller files.
type Controller struct {
	*quimby.Flagger
	name  string
	stubs bool
}

func (c *Controller) Desc() string {
	return "Create a new controller. Camel-cased names will be used unmodified. Other names will be title-cased."
}

func (c *Controller) SetFlags() {
	c.StringVar(&c.name, "name", "", "Name of the controller type to create")
	c.BoolVar(&c.stubs, "stubs", false, "Decides whether the generated controller will have stubs for the 5 standard action methods")
}

func (c *Controller) Run() {
	name := c.name
	if name == "" {
		args := c.FlagSet.Args()
		if len(args) == 0 {
			fmt.Println("You must supply a name for the controller.")
			return
		}
		name = args[0]
	}
	current, _ := os.Getwd()
	gopath := filepath.Join(os.ExpandEnv("$GOPATH"), "src")
	if !strings.HasPrefix(current, gopath) {
		fmt.Println("Controllers must be created in a Gadget project directory.")
		return
	}
	createControllerFile(current, name, c.stubs)
}

func createControllerFile(projectName, controllerName string, stubs bool) {
	var filename string
	t := getTemplate("controller.tpl")
	if strutil.IsPascalCase(controllerName) {
		filename = strutil.Snakify(controllerName)
	} else {
		filename = strings.ToLower(controllerName)
		controllerName = strings.Title(controllerName)
	}
	path := filepath.Join(projectName, "controllers", filename+".go")
	importPath, _ := getImportPath(projectName)
	if f, err := os.Create(path); err != nil {
		fmt.Printf("Unable to create file controllers/%s.go: %s\n", filename, err)
	} else {
		defer f.Close()
		t.Execute(f, struct{
			Name, Project string
			Stubs bool
		}{
			Name:    controllerName,
			Project: importPath,
			Stubs:   stubs,
		})
		fmt.Printf("Created %s\n", path)
	}
}
