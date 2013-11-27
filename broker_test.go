package gadget

import (
	. "launchpad.net/gocheck"
)

type ProcessSuite struct{}

type process struct {
	*App
}

var p *process

func (s *ProcessSuite) SetUpTest(c *C) {
	p = &process{&App{}}
}

func (s *ProcessSuite) TearDownTest(c *C) {
	p.Brokers = make(map[string]Broker)
}

var _ = Suite(&ProcessSuite{})

//Calling `Process(200, "<html></html>", "text/html")` should return 200, "<html></html>", false if there is no HTML processor
func (s *ProcessSuite) TestCallingProcesstexthtml200HtmlhtmlShouldReturn200HtmlhtmlFalseThereIsNoHtmlProcessor(c *C) {
	status, body, matched, changed := p.process(&Request{}, 200, "<html></html>", "text/html", &RouteData{})
	c.Assert(status, Equals, 200)
	c.Assert(body, Equals, "<html></html>")
	c.Assert(matched, Equals, "text/html; charset=utf-8")
	c.Assert(changed, Equals, false)
}

//Calling `Process(200, []string{"foo", "bar", "baz"}, "application/json")` should return 200, "[foo bar baz]", false if there is no JSON processor
func (s *ProcessSuite) TestCallingProcessapplicationjson200StringfooBarBazShouldReturn200FooBarBazFalseThereIsNoJsonBroker(c *C) {
	status, body, matched, changed := p.process(&Request{}, 200, []string{"foo", "bar", "baz"}, "application/json", &RouteData{})
	c.Assert(status, Equals, 200)
	c.Assert(body, Equals, "[foo bar baz]")
	c.Assert(matched, Equals, "text/plain; charset=utf-8")
	c.Assert(changed, Equals, false)
}

//Calling `Process(200, "hi there", "application/json")` should encode the string as a JSON value if there is a JSON processor
func (s *ProcessSuite) TestCallingProcessapplicationjson200HiThereShouldReturn200HiThereFalseThereIsNoJsonBroker(c *C) {
	p.Accept("application/json").Via(JsonBroker)
	status, body, matched, changed := p.process(&Request{}, 200, "hi there", "application/json", &RouteData{})
	c.Assert(status, Equals, 200)
	c.Assert(body, Equals, `"hi there"`)
	c.Assert(matched, Equals, "application/json")
	c.Assert(changed, Equals, true)
}

//Calling `Process("application/json", 200, []string{"foo", "bar", "baz"})` should return 200, '["foo", "bar", "baz"]', true when there is a JSON processor
func (s *ProcessSuite) TestCallingProcessapplicationjson200StringfooBarBazShouldReturn200FooBarBazTrueWhenThereIsJsonBroker(c *C) {
	p.Accept("application/json").Via(JsonBroker)
	status, body, matched, changed := p.process(&Request{}, 200, []string{"foo", "bar", "baz"}, "application/json", &RouteData{})
	c.Assert(status, Equals, 200)
	c.Assert(body, Equals, `["foo","bar","baz"]`)
	c.Assert(matched, Equals, "application/json")
	c.Assert(changed, Equals, true)
}

//Calling `Process` with a function for the body value should return 500, "", true
func (s *ProcessSuite) TestCallingProcessAnonymousFunctionForBodyValueShouldReturn500True(c *C) {
	p.Accept("application/json").Via(JsonBroker)
	status, body, matched, changed := p.process(&Request{}, 200, JsonBroker, "application/json", &RouteData{})
	c.Assert(status, Equals, 500)
	c.Assert(body, Equals, "")
	c.Assert(matched, Equals, "application/json")
	c.Assert(changed, Equals, true)
}
