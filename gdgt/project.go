package main

import (
	"fmt"
	"github.com/redneckbeard/quimby"
	"os"
	"path/filepath"
)

func init() {
	quimby.Add(&New{})
}

type New struct {
	*quimby.Flagger
	name string
}

func (c *New) Desc() string {
	return "Create a new Gadget project."
}

func (c *New) SetFlags() {
	c.StringVar(&c.name, "name", "", "Name of the project to be created")
}

func (c *New) Run() {
	name := c.name
	if name == "" {
		args := c.FlagSet.Args()
		if len(args) == 0 {
			fmt.Println("You must supply a name for your new project.")
			return
		}
		name = args[0]
	}
	current, _ := os.Getwd()
	path, err := filepath.Rel(
		filepath.Join(os.ExpandEnv("$GOPATH"), "src"),
		filepath.Join(current, name),
	)
	if err != nil {
		fmt.Println("Projects must be created in the src directory of your GOPATH.")
	}
	subdirs := []string{
		"controllers",
		"app",
		"templates/home",
		"static/css",
		"static/img",
		"static/js",
	}
	for _, subdir := range subdirs {
		path := filepath.Join(name, subdir)
		err := os.MkdirAll(path, 0777)
		if err != nil {
			fmt.Printf("Could not create directory %s: %s\n", path, err)
		} else {
			fmt.Printf("Created directory %s\n", path)
		}
	}
	t := getTemplate("main.tpl")
	if f, err := os.Create(filepath.Join(name, "main.go")); err != nil {
		fmt.Printf("Unable to create file main.go in %s: %s\n", name, err)
	} else {
		defer f.Close()
		t.Execute(f, map[string]string{
			"projectName": name,
			"path": path,
		})
		fmt.Printf("Created %s/main.go\n", name)
	}
	copyTemplate("conf.tpl", filepath.Join(name, "app", "conf.go"))
	copyTemplate("base.html", filepath.Join(name, "templates", "base.html"))
	copyTemplate("home/index.html", filepath.Join(name, "templates", "home", "index.html"))
	fmt.Println("Created app/conf.go")
	createControllerFile(name, "Home")
}


