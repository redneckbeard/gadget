package templates

import (
	"github.com/redneckbeard/gadget"
	. "launchpad.net/gocheck"
	"strings"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type TemplateSuite struct{}

var _ = Suite(&TemplateSuite{})

//A 20x status code should be returned for requests to endpoints that have existing templates
func (s *TemplateSuite) Test20x(c *C) {
	TemplatePath = "testdata/200"
	for i := 200; i < 300; i++ {
		status, _ := TemplateBroker(&gadget.Request{}, i, "body", &gadget.RouteData{"widgets", "index", "GET"})
		c.Assert(status, Equals, i)
	}
}

// Calling "render" in a template should load a subtemplate from the templates/$controllerName directory, falling back to the templates directory if it does not exist
func (s *TemplateSuite) TestRenderHelper(c *C) {
	templateDirs := []string{
		"testdata/render/controllerdir",
		"testdata/render/fallback",
	}
	for _, dir := range templateDirs {
		TemplatePath = dir
		context := struct{
			Message string
		}{
			Message: "context passed to subtemplate",
		}
		status, body := TemplateBroker(&gadget.Request{}, 200, context, &gadget.RouteData{"widgets", "index", "GET"})
		c.Assert(status, Equals, 200)
		c.Assert(strings.TrimSpace(body), Equals, context.Message)
	}
}

// Templates are successfully rendered if they do not define templates invoked by the base template
func (s *TemplateSuite) TestMissingDefine(c *C) {
	TemplatePath = "testdata/define"
	context := "simple"
	status, body := TemplateBroker(&gadget.Request{}, 200, context, &gadget.RouteData{"widgets", "index", "GET"})
	c.Assert(status, Equals, 200)
	c.Assert(strings.TrimSpace(body), Equals, context)
}
