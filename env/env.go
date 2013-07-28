package env

import (
	"flag"
	"github.com/redneckbeard/gadget/cmd"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var (
	root, staticPrefix, logFilePath, Port string
	logger                                *log.Logger
	Debug                                 bool
	Handler                               http.HandlerFunc
	configured                            bool
)

func init() {
	cmd.Add("serve", &Serve{FlagSet: &flag.FlagSet{}})
}

type Serve struct{
	*flag.FlagSet
}

func (s *Serve) SetFlags() {
	s.StringVar(&staticPrefix, "static", "/static/", "URL prefix for serving the 'static' directory")
	s.StringVar(&root, "root", "", "Directory that contains uncompiled application assets")
	s.StringVar(&logFilePath, "log", "", "Path to log file")
	s.StringVar(&Port, "port", "8090", "Port on which the application will listen")
	s.BoolVar(&Debug, "debug", true, "Sets the env.Debug value within Gadget")
}

func (s *Serve) Run() {
	if root == "" {
		if wd, err := os.Getwd(); err != nil {
			panic(err)
		} else {
			root = wd
		}
	} else if !filepath.IsAbs(root) {
		panic("fileroot must be an absolute path")
	}
	var writer io.Writer
	if logFilePath != "" {
		if !filepath.IsAbs(logFilePath) {
			logFilePath = AbsPath(logFilePath)
		}
		if f, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE, 0666); err != nil {
			panic(err)
		} else {
			writer = f
		}
	} else {
		writer = os.Stdout
	}
	logger = log.New(writer, "gdgt| ", 0)
	serveStatic()
	http.HandleFunc("/", Handler)
	err := http.ListenAndServe(":"+Port, nil)
	if err != nil {
		panic(err)
	}
}

func AbsPath(path ...string) string {
	return filepath.Join(append([]string{root}, path...)...)
}

func serveStatic() {
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
