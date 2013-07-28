package templates

import (
	"bytes"
	"github.com/redneckbeard/gadget"
	"github.com/redneckbeard/gadget/env"
	"html/template"
	"strconv"
)

var registry = make(template.FuncMap)

func AddHelper(name string, f interface{}) {
	registry[name] = f
}

func templatePath(components ...string) string {
	components[len(components)-1] += ".html"
	return env.AbsPath(append([]string{"templates"}, components...)...)
}

func TemplateBroker(r *gadget.Request, status int, body interface{}, data *gadget.RouteData) (int, string) {
	helpers := template.FuncMap{
		"request": func() *gadget.Request { return r },
	}
	for name, helper := range registry {
		helpers[name] = helper
	}
	t, err := template.New("base.html").Funcs(helpers).ParseFiles(templatePath("base"))
	if err != nil {
		return 404, err.Error()
	}
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
