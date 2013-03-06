package env

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

var (
	e                  = &Env{}
	root, staticPrefix string
)

type Env struct {
	FileRoot, StaticPrefix string
}

func Configure() error {
	flag.StringVar(&staticPrefix, "static", "/static/", "URL prefix for serving the 'static' directory")
	flag.StringVar(&root, "root", "", "Directory that contains uncompiled application assets")
	flag.Parse()
	if root == "" {
		root, err := os.Getwd()
		if err != nil {
			return err
		}
		fmt.Println(root)
		e.FileRoot = root
	} else if !filepath.IsAbs(root) {
		return errors.New("fileroot must be an absolute path")
	} else {
		e.FileRoot = root
	}
	e.StaticPrefix = staticPrefix
	return nil
}

func AbsPath(path ...string) string {
	return filepath.Join(append([]string{e.FileRoot}, path...)...)
}

func ServeStatic() {
	http.Handle(e.StaticPrefix, http.StripPrefix(e.StaticPrefix, http.FileServer(http.Dir(AbsPath("static")))))
}

func Open(path string) (*os.File, error) {
	return os.Open(AbsPath(path))
}
