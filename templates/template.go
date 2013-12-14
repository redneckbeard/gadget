package templates

import (
	"bytes"
	"github.com/redneckbeard/gadget"
	"github.com/redneckbeard/gadget/env"
	"html/template"
	"strconv"
)

var (
	registry     = make(template.FuncMap)
	TemplatePath = "templates"
)

// AddHelper registers functions with TemplateBroker that will be available in
// templates during rendering.
func AddHelper(name string, f interface{}) {
	registry[name] = f
}

func templatePath(components ...string) string {
	components[len(components)-1] += ".html"
	return env.RelPath(append([]string{TemplatePath}, components...)...)
}

// TemplateBroker attempts to render interface{} value body as the context of a
// html/template.Template. It requires adherence to a few simple conventions
// for locating templates: 1) all templates are inside a "templates" directory
// in the root directory of the Gadget application; 2) a "templates/base.html"
// file exists and can be parsed as a Go template; 3) templates will be loaded
// from a directory matching the lower-cased plural name of the controller and
// the lower-cased name of the action, plus the extension ".html". For example,
// the Show method of a FavoriteController would look for a template
// "templates/favorites/show.html."
//
// All error codes can also be served via their own templates. Non-200 statuses
// will result in TemplateBroker looking for a "templates/403.html",
// "templates/502.html", etc.
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
	if status >= 200 && status < 300 {
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
