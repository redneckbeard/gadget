package routing

import (
	"github.com/redneckbeard/gadget/controller"
	"github.com/redneckbeard/gadget/processor"
	"github.com/redneckbeard/gadget/requests"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"net/http"
	"net/http/httptest"
)

type HandlerSuite struct{}

func (s *HandlerSuite) SetUpSuite(c *C) {
	controller.Register(&MapController{controller.New()})
	processor.Define("application/json", processor.JsonProcessor)
	Routes(SetIndex("map"))
}

var _ = Suite(&HandlerSuite{})

type MapController struct {
	*controller.DefaultController
}

func (c *MapController) Index(r *requests.Request) (int, interface{}) {
	retVal := make(map[string]int)
	retVal["foo"] = 1
	retVal["bar"] = 2
	return 200, retVal
}

//A response should be run through the JSON processor when one is defined and when the request is made with Accepts: application/json
func (s *HandlerSuite) TestResponseRunThroughJsonProcessorWhenOneIsDefinedAndWhenRequestIsMadeAcceptsApplicationjson(c *C) {
	handler := Handler()

	req, err := http.NewRequest("GET", "http://127.0.0.1:8000/", nil)
	c.Assert(err, IsNil)
	req.Header.Set("Accept", "application/json")

	resp := httptest.NewRecorder()
	handler(resp, req)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, `{"bar":2,"foo":1}`)
	c.Assert(resp.Header().Get("Content-Type"), Equals, "application/json")
}

//A response should be run through the JSON processor when one is defined and when the request is made with Content-Type: application/json
func (s *HandlerSuite) TestResponseRunThroughJsonProcessorWhenOneIsDefinedAndWhenRequestIsMadeContenttypeApplicationjson(c *C) {
	handler := Handler()

	req, err := http.NewRequest("GET", "http://127.0.0.1:8000/", nil)
	c.Assert(err, IsNil)
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	handler(resp, req)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, `{"bar":2,"foo":1}`)
	c.Assert(resp.Header().Get("Content-Type"), Equals, "application/json")
}

//A response should be run through the JSON processor when one is defined and when the request is made with Accepts: application/json and Content-Type: text/xml
func (s *HandlerSuite) TestResponseRunThroughJsonProcessorWhenOneIsDefinedAndWhenRequestIsMadeAcceptsApplicationjsonAndContenttypeTextxml(c *C) {
	handler := Handler()

	req, err := http.NewRequest("GET", "http://127.0.0.1:8000/", nil)
	c.Assert(err, IsNil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "text/xml")

	resp := httptest.NewRecorder()
	handler(resp, req)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, `{"bar":2,"foo":1}`)
	c.Assert(resp.Header().Get("Content-Type"), Equals, "application/json")
}

//A response should not be run through the JSON processor when one is defined and when the request is made with Content-Type: application/json and Accepts: text/xml
func (s *HandlerSuite) TestResponseRunThroughJsonProcessorWhenOneIsDefinedAndWhenRequestIsMadeContenttypeApplicationjsonAndAcceptsTextxml(c *C) {
	handler := Handler()

	req, err := http.NewRequest("GET", "http://127.0.0.1:8000/", nil)
	c.Assert(err, IsNil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/xml")

	resp := httptest.NewRecorder()
	handler(resp, req)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "map[bar:2 foo:1]")
	c.Assert(resp.Header().Get("Content-Type"), Equals, "text/html")
}
