package routing

import (
	"fmt"
	"github.com/redneckbeard/gadget/controller"
	"github.com/redneckbeard/gadget/processor"
	"github.com/redneckbeard/gadget/requests"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"net/http"
	"net/http/httptest"
	"crypto/md5"
	"strings"
)

type HandlerSuite struct{}

func (s *HandlerSuite) SetUpSuite(c *C) {
	controller.Register(&MapController{controller.New()})
	controller.Register(&ResourceController{controller.New()})
	controller.Register(&UuidController{controller.New()})
	processor.Define("application/json", processor.JsonProcessor)
	Routes(SetIndex("map"), Resource("resource"), Resource("uuid"))
}

var _ = Suite(&HandlerSuite{})

type ResourceController struct{ *controller.DefaultController }

func (c *ResourceController) Index(r *requests.Request) (int, interface{}) { return 200, "" }
func (c *ResourceController) Show(r *requests.Request) (int, interface{})  { return 200, "" }
func (c *ResourceController) Extra(r *requests.Request) (int, interface{}) { return 200, "" }

type UuidController struct{ *controller.DefaultController }

func (c *UuidController) IdPattern() string { return `\w{8}-\w{4}-\w{4}-\w{4}-\w{12}` }

func (c *UuidController) Index(r *requests.Request) (int, interface{}) { return 200, "" }
func (c *UuidController) Show(r *requests.Request) (int, interface{})  { return 200, "" }
func (c *UuidController) Extra(r *requests.Request) (int, interface{}) { return 200, "" }



type MapController struct {
	*controller.DefaultController
}

func (c *MapController) Index(r *requests.Request) (int, interface{}) {
	retVal := &struct {
		Bar int `json:"bar"`
		Foo int `json:"foo"`
	}{2, 1}
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
	c.Assert(string(body), Equals, "&{2 1}")
	c.Assert(resp.Header().Get("Content-Type"), Equals, "text/plain; charset=utf-8")
}

//Route.Respond should 404 on a component that is neither an ID match nor an action match
func (s *HandlerSuite) TestRoute404sNoIdNoAction(c *C) {
	handler := Handler()

	req, _ := http.NewRequest("GET", "http://127.0.0.1:8000/resource/not-extra", nil)
	resp := httptest.NewRecorder()
	handler(resp, req)
	c.Assert(resp.Code, Equals, 404)
}

func fakeUuid() string {
	h := md5.New()
	sum := fmt.Sprintf("%x", h.Sum(nil))
	return strings.Join([]string{sum[:8], sum[8:12], sum[12:16], sum[16:20], sum[20:]}, "-")
}

func (s *HandlerSuite) TestCustomIdPattern200OnMatches(c *C) {
	handler := Handler()

	req, _ := http.NewRequest("GET", "http://127.0.0.1:8000/uuid/" + fakeUuid(), nil)
	resp := httptest.NewRecorder()
	handler(resp, req)
	c.Assert(resp.Code, Equals, 200)
}

func (s *HandlerSuite) TestCustomIdPattern404sOnNonMatches(c *C) {
	handler := Handler()

	req, _ := http.NewRequest("GET", "http://127.0.0.1:8000/uuid/" + fakeUuid()[1:], nil)
	resp := httptest.NewRecorder()
	handler(resp, req)
	c.Assert(resp.Code, Equals, 404)
}

