package gadget

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"net/http"
	"net/http/httptest"
	"strings"
)

type HandlerSuite struct{}

type handlerApp struct {
	*App
}

var h *handlerApp

func (s *HandlerSuite) SetUpSuite(c *C) {
	h = &handlerApp{ &App{} }
	h.Register(&MapController{})
	h.Register(&ResourceController{})
	h.Register(&UuidController{})
	h.Accept("application/json").Via(JsonBroker)
	h.Routes(h.SetIndex("maps"), h.Resource("resources"), h.Resource("uuids"))
}
func (s *HandlerSuite) TearDownSuite(c *C) {
	h.Controllers = make(map[string]Controller)
}

var _ = Suite(&HandlerSuite{})

type ResourceController struct{ *DefaultController }

func (c *ResourceController) Index(r *Request) (int, interface{})      { return 200, "" }
func (c *ResourceController) Show(r *Request) (int, interface{})       { return 200, "" }
func (c *ResourceController) Extra(r *Request) (int, interface{})      { return 200, "" }
func (c *ResourceController) PascalCase(r *Request) (int, interface{}) { return 200, "" }

type UuidController struct{ *DefaultController }

func (c *UuidController) IdPattern() string { return `\w{8}-\w{4}-\w{4}-\w{4}-\w{12}` }

func (c *UuidController) Index(r *Request) (int, interface{}) { return 200, "" }
func (c *UuidController) Show(r *Request) (int, interface{})  { return 200, "" }
func (c *UuidController) Extra(r *Request) (int, interface{}) { return 200, "" }

type MapController struct {
	*DefaultController
}

func (c *MapController) Index(r *Request) (int, interface{}) {
	retVal := &struct {
		Bar int `json:"bar"`
		Foo int `json:"foo"`
	}{2, 1}
	return 200, retVal
}

//A response should be run through the JSON processor when one is defined and when the request is made with Accepts: application/json
func (s *HandlerSuite) TestResponseRunThroughJsonProcessorWhenOneIsDefinedAndWhenRequestIsMadeAcceptsApplicationjson(c *C) {
	handler := h.Handler()

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
	handler := h.Handler()

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
	handler := h.Handler()

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
	handler := h.Handler()

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
	handler := h.Handler()

	req, _ := http.NewRequest("GET", "http://127.0.0.1:8000/resources/not-extra", nil)
	resp := httptest.NewRecorder()
	handler(resp, req)
	c.Assert(resp.Code, Equals, 404)
}

//Route.Respond should 404 on a component that matches a default action name
func (s *HandlerSuite) TestRouterespondShould404OnComponentThatMatchesDefaultActionName(c *C) {
	handler := h.Handler()

	req, _ := http.NewRequest("GET", "http://127.0.0.1:8000/resources/index", nil)
	resp := httptest.NewRecorder()
	handler(resp, req)
	c.Assert(resp.Code, Equals, 404)
}

//Hyphenated URL components should be correctly routed to Pascal-cased method names
func (s *HandlerSuite) TestHyphenatedUrlComponentsCorrectlyRoutedToPascalcasedMethodNames(c *C) {
	handler := h.Handler()

	req, _ := http.NewRequest("GET", "http://127.0.0.1:8000/resources/pascal-case", nil)
	resp := httptest.NewRecorder()
	handler(resp, req)
	c.Assert(resp.Code, Equals, 200)
}

//Route.Respond should 404 on a component that matches an exported method but does not have an Action signature
func (s *HandlerSuite) TestRouterespondShould404OnComponentThatMatchesExportedMethodButDoesNotHaveActionSignature(c *C) {
	handler := h.Handler()

	req, _ := http.NewRequest("GET", "http://127.0.0.1:8000/resources/idpattern", nil)
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
	handler := h.Handler()

	req, _ := http.NewRequest("GET", "http://127.0.0.1:8000/uuids/"+fakeUuid(), nil)
	resp := httptest.NewRecorder()
	handler(resp, req)
	c.Assert(resp.Code, Equals, 200)
}

func (s *HandlerSuite) TestCustomIdPattern404sOnNonMatches(c *C) {
	handler := h.Handler()

	req, _ := http.NewRequest("GET", "http://127.0.0.1:8000/uuids/"+fakeUuid()[1:], nil)
	resp := httptest.NewRecorder()
	handler(resp, req)
	c.Assert(resp.Code, Equals, 404)
}
