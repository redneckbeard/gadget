package processor

import (
	. "launchpad.net/gocheck"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type ProcessSuite struct{}

func (s *ProcessSuite) TearDownTest(c *C) {
	clear()
}

var _ = Suite(&ProcessSuite{})

//Calling `Process(200, "<html></html>", "text/html")` should return 200, "<html></html>", false if there is no HTML processor
func (s *ProcessSuite) TestCallingProcesstexthtml200HtmlhtmlShouldReturn200HtmlhtmlFalseThereIsNoHtmlProcessor(c *C) {
	status, body, changed := Process(200, "<html></html>", "text/html")
	c.Assert(status, Equals, 200)
	c.Assert(body, Equals, "<html></html>")
	c.Assert(changed, Equals, false)
}

//Calling `Process(200, []string{"foo", "bar", "baz"}, "application/json")` should return 200, "[foo bar baz]", false if there is no JSON processor
func (s *ProcessSuite) TestCallingProcessapplicationjson200StringfooBarBazShouldReturn200FooBarBazFalseThereIsNoJsonProcessor(c *C) {
	status, body, changed := Process(200, []string{"foo", "bar", "baz"}, "application/json")
	c.Assert(status, Equals, 200)
	c.Assert(body, Equals, "[foo bar baz]")
	c.Assert(changed, Equals, false)
}

//Calling `Process(200, "hi there", "application/json")` should return 200, "hi there", false if there is a JSON processor
func (s *ProcessSuite) TestCallingProcessapplicationjson200HiThereShouldReturn200HiThereFalseThereIsNoJsonProcessor(c *C) {
	Define("application/json", JsonProcessor)
	status, body, changed := Process(200, "hi there", "application/json")
	c.Assert(status, Equals, 200)
	c.Assert(body, Equals, "hi there")
	c.Assert(changed, Equals, false)
}

//Calling `Process("application/json", 200, []string{"foo", "bar", "baz"})` should return 200, '["foo", "bar", "baz"]', true when there is a JSON processor
func (s *ProcessSuite) TestCallingProcessapplicationjson200StringfooBarBazShouldReturn200FooBarBazTrueWhenThereIsJsonProcessor(c *C) {
	Define("application/json", JsonProcessor)
	status, body, changed := Process(200, []string{"foo", "bar", "baz"}, "application/json")
	c.Assert(status, Equals, 200)
	c.Assert(body, Equals, `["foo","bar","baz"]`)
	c.Assert(changed, Equals, true)
}

//Calling `Process` with a function for the body value should return 500, "", true
func (s *ProcessSuite) TestCallingProcessAnonymousFunctionForBodyValueShouldReturn500True(c *C){
	Define("application/json", JsonProcessor)
	status, body, changed := Process(200, JsonProcessor, "application/json")
	c.Assert(status, Equals, 500)
	c.Assert(body, Equals, "")
	c.Assert(changed, Equals, true)
}


