package templates

import (
	. "launchpad.net/gocheck"
	"github.com/redneckbeard/gadget"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type TemplateSuite struct{}

var _ = Suite(&TemplateSuite{})

//A 20x status code should be returned for requests to endpoints that have existing templates
func (s *TemplateSuite) Test20x(c *C) {
	for i := 200; i < 300; i++ {
		status, _ := TemplateBroker(&gadget.Request{}, i, "body", &gadget.RouteData{"widgets", "index", "GET"})
		c.Assert(status, Equals, i)
	}
}
