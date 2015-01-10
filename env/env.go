package env

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/redneckbeard/quimby"
)

var (
	root, staticPrefix, logFilePath, port string
	logger                                *log.Logger
	messages                              = make(chan []interface{})
	envVars                               map[string]string
	// Debug is set via the -debug flag for the serve command.
	Debug bool
	// Handler comes from calling Handler() on a gadget.App object. It's used by the serve command to run the server.
	Handler    http.HandlerFunc
	configured bool
)

func init() {
	quimby.Add(&Serve{})
}

// The Serve command makes it easy to run Gadget applications.
type Serve struct {
	*quimby.Flagger
}

func (s *Serve) Desc() string {
	return "Start a gadget server."
}

// SetFlags defines flags for the serve command.
func (s *Serve) SetFlags() {
	s.StringVar(&staticPrefix, "static", "/static/", "URL prefix for serving the 'static' directory")
	s.StringVar(&root, "root", "", "Directory that contains uncompiled application assets. Defaults to current working directory.")
	s.StringVar(&logFilePath, "log", "", "Path to log file")
	s.StringVar(&port, "port", "8090", "port on which the application will listen")
	s.BoolVar(&Debug, "debug", true, "Sets the env.Debug value within Gadget")
}

// Run sets up a logger and runs the Handler.
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
			logFilePath = RelPath(logFilePath)
		}
		if f, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE, 0666); err != nil {
			panic(err)
		} else {
			writer = f
		}
	} else {
		writer = os.Stdout
	}
	logger = log.New(writer, "", 0)
	go func() {
		for msg := range messages {
			logger.Println(msg...)
		}
	}()
	serveStatic()
	http.HandleFunc("/", Handler)
	Log("Running Gadget at 0.0.0.0:" + port + "...")
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}

// RelPath creates an absolute path from path segments path relative to the project root.
func RelPath(path ...string) string {
	return filepath.Join(append([]string{root}, path...)...)
}

func serveStatic() {
	http.Handle(staticPrefix, http.StripPrefix(staticPrefix, http.FileServer(http.Dir(RelPath("static")))))
}

// Open wraps os.Open, but with the assumption that the path is relative to the project root.
func Open(path string) (*os.File, error) {
	return os.Open(RelPath(path))
}

// Log writes arguments v as a single line to the default logger.
func Log(v ...interface{}) {
	go func() { messages <- v }()
}

func Get(varname string) string {
	if Debug {
		if envVars == nil {
			envVars = make(map[string]string)
			f, err := Open(".env")
			defer f.Close()
			if err != nil {
				return ""
			}
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				text := scanner.Text()
				pair := strings.SplitN(text, "=", 2)
				envVars[pair[0]] = pair[1]
			}
		}
		return envVars[varname]
	}
	return os.Getenv(varname)
}
