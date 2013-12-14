package templates

import (
	. "launchpad.net/gocheck"
	"github.com/redneckbeard/gadget"
	"testing"
	"strings"
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
		context := "context passed to subtemplate"
		status, body := TemplateBroker(&gadget.Request{}, 200, context, &gadget.RouteData{"widgets", "index", "GET"})
		c.Assert(status, Equals, 200)
		c.Assert(strings.TrimSpace(body), Equals, context)
	}
}
