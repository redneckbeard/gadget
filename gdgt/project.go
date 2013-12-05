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

func getGoPath(name string) (string, error) {
	var path string
	gopath := os.ExpandEnv("$GOPATH")
	if gopath == "" {
		return "", fmt.Errorf("GOPATH must be set to use the gdgt tool.")
	}
	inGopath := filepath.Join(gopath, "src", "*")
	current, _ := os.Getwd()
	if matched, _ := filepath.Match(inGopath, current); matched {
		path = filepath.Join(current, name)
	} else {
		path = filepath.Join(gopath, "src", name)
	}
	return path, nil
}

func getImportPath(name string) (string, error) {
	return filepath.Rel(filepath.Join(os.ExpandEnv("$GOPATH"), "src"), name)
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

	path, err := getGoPath(name) 
	if err != nil {
		fmt.Println(err)
		return
	}

	importPath, _ := getImportPath(path)
	subdirs := []string{
		"controllers",
		"app",
		"templates/home",
		"static/css",
		"static/img",
		"static/js",
	}
	for _, subdir := range subdirs {
		path := filepath.Join(path, subdir)
		err := os.MkdirAll(path, 0777)
		if err != nil {
			fmt.Printf("Could not create directory %s: %s\n", path, err)
		} else {
			fmt.Printf("Created directory %s\n", path)
		}
	}
	t := getTemplate("main.tpl")
	if f, err := os.Create(filepath.Join(path, "main.go")); err != nil {
		fmt.Printf("Unable to create file main.go in %s: %s\n", name, err)
	} else {
		defer f.Close()
		t.Execute(f, map[string]string{
			"projectName": name,
			"path": importPath,
		})
		fmt.Printf("Created %s/main.go\n", path)
	}
	copyTemplate("conf.tpl", filepath.Join(path, "app", "conf.go"))
	copyTemplate("base.html", filepath.Join(path, "templates", "base.html"))
	copyTemplate("404.html", filepath.Join(path, "templates", "404.html"))
	copyTemplate("home/index.html", filepath.Join(path, "templates", "home", "index.html"))
	fmt.Printf("Created %s/app/conf.go\n", path)
	createControllerFile(path, "Home", false)
}
