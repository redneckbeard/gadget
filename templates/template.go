package templates

import (
	"bytes"
	"github.com/redneckbeard/gadget/env"
	"github.com/redneckbeard/gadget/processor"
	"html/template"
	"strconv"
)

var helpers = make(template.FuncMap)

func AddHelper(name string, f interface{}) {
	helpers[name] = f
}

func templatePath(components ...string) string {
	components[len(components)-1] += ".html"
	return env.AbsPath(append([]string{"templates"}, components...)...)
}

func TemplateProcessor(status int, body interface{}, data *processor.RouteData) (int, string) {
	t, err := template.ParseFiles(templatePath("base"))
	if err != nil {
		return 404, err.Error()
	}
	t = t.Funcs(helpers)
	var mainTemplatePath string
	if status == 200 {
		mainTemplatePath = templatePath(data.ControllerName, data.Action)
	} else {
		mainTemplatePath = templatePath(strconv.FormatInt(int64(status), 10))
	}
	_, err = t.ParseFiles(mainTemplatePath)
	if err != nil {
		return 500, err.Error()
	}
	buf := new(bytes.Buffer)
	err = t.Execute(buf, body)
	if err != nil {
		return 500, err.Error()
	}
	return status, string(buf.Bytes())
}
