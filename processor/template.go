package processor

import (
	"bytes"
	"github.com/redneckbeard/gadget/env"
	"html/template"
)

func templatePath(components ...string) string {
	components[len(components)-1] += ".html"
	return env.AbsPath(append([]string{"templates"}, components...)...)
}

func TemplateProcessor(status int, body interface{}, data *RouteData) (int, string) {
	t, err := template.ParseFiles(templatePath("base"))
	if err != nil {
		return 404, err.Error()
	}
	_, err = t.ParseFiles(templatePath(data.ControllerName, data.Action))
	if err != nil {
		return 404, err.Error()
	}
	buf := new(bytes.Buffer)
	err = t.Execute(buf, body)
	if err != nil {
		return 500, err.Error()
	}
	return 200, string(buf.Bytes())
}
