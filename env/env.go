package env

import (
	"errors"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var (
	root, staticPrefix, logFilePath string
	logger                          *log.Logger
)

func Configure() error {
	flag.StringVar(&staticPrefix, "static", "/static/", "URL prefix for serving the 'static' directory")
	flag.StringVar(&root, "root", "", "Directory that contains uncompiled application assets")
	flag.StringVar(&logFilePath, "log", "", "Path to log file")
	flag.Parse()
	if root == "" {
		if wd, err := os.Getwd(); err != nil {
			return err
		} else {
			root = wd
		}
	} else if !filepath.IsAbs(root) {
		return errors.New("fileroot must be an absolute path")
	}
	var writer io.Writer
	if logFilePath != "" {
		if !filepath.IsAbs(logFilePath) {
			logFilePath = AbsPath(logFilePath)
		} 
		if f, err := os.OpenFile(logFilePath, os.O_RDWR | os.O_CREATE, 0666); err != nil {
			return err
		} else {
			writer = f
		}
	} else {
		writer = os.Stdout
	}
	logger = log.New(writer, "gdgt| ", 0)
	return nil
}

func AbsPath(path ...string) string {
	return filepath.Join(append([]string{root}, path...)...)
}

func ServeStatic() {
	http.Handle(staticPrefix, http.StripPrefix(staticPrefix, http.FileServer(http.Dir(AbsPath("static")))))
}

func Open(path string) (*os.File, error) {
	return os.Open(AbsPath(path))
}

func Log(v ...interface{}) {
	if logger != nil {
		logger.Println(v...)
	}
}
