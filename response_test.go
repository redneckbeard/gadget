package gadget

import (
	"github.com/redneckbeard/gadget/processor"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"net/http"
	"net/http/httptest"
	"time"
)

type ResponseSuite struct{}

var _ = Suite(&ResponseSuite{})

func (s *ResponseSuite) SetUpSuite(c *C) {
	Register(&ResponseController{})
	Register(&ImplicitController{})
	processor.Define("application/json", processor.JsonProcessor)
	Routes(Resource("responses"), Resource("implicits"))
}
func (s *ResponseSuite) TearDownSuite(c *C) {
	clear()
}

var cookie = &http.Cookie{
	Name:    "foo",
	Value:   "bar",
	Expires: time.Now().Add(time.Duration(10 * time.Hour)),
}

type ResponseController struct {
	*DefaultController
}

func (c *ResponseController) Index(*Request) (int, interface{}) {
	body := struct{ Foo, Bar string }{"baz", "quux"}
	response := NewResponse(body)
	response.Headers.Set("X-Framework", "Gadget")
	return 200, response
}

func (c *ResponseController) Show(*Request) (int, interface{}) {
	response := NewResponse("test")
	response.AddCookie(cookie)
	return 200, response
}

type ImplicitController struct {
	*DefaultController
}

func (c *ImplicitController) Index(*Request) (int, interface{}) {
	body := struct{ Foo, Bar string }{"baz", "quux"}
	return 200, body
}

//Headers set on a gadget.Response in a controller method are correctly transferred to the http.Response
func (s *ResponseSuite) TestHeadersSetOnGadgetresponseControllerMethodAreCorrectlyTransferredToHttpresponse(c *C) {
	handler := Handler()

	req, err := http.NewRequest("GET", "http://127.0.0.1:8000/responses", nil)
	c.Assert(err, IsNil)

	resp := httptest.NewRecorder()
	handler(resp, req)
	c.Assert(resp.Header().Get("X-Framework"), Equals, "Gadget")
}

//Cookies added to a gadget.Response in a controller method are correctly transferred to the http.Response
func (s *ResponseSuite) TestCookiesAddedToGadgetresponseControllerMethodAreCorrectlyTransferredToHttpresponse(c *C) {
	handler := Handler()

	req, err := http.NewRequest("GET", "http://127.0.0.1:8000/responses/1", nil)
	c.Assert(err, IsNil)

	resp := httptest.NewRecorder()
	handler(resp, req)
	c.Assert(resp.Code, Equals, 200)
	c.Assert(resp.Header().Get("Set-Cookie"), Equals, cookie.String())
}

//The body of the http.Response is identical between a controller method that returns an anonymous struct and one that returns a gadget.Response with its Body set to that struct
func (s *ResponseSuite) TestBodyHttpresponseIsIdenticalBetweenControllerMethodThatReturnsAnonymousStructAndOneThatReturnsGadgetresponseItsBodySetToThatStruct(c *C) {
	handler := Handler()

	req, err := http.NewRequest("GET", "http://127.0.0.1:8000/responses", nil)
	c.Assert(err, IsNil)
	resp := httptest.NewRecorder()
	handler(resp, req)
	body1, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)

	req, err = http.NewRequest("GET", "http://127.0.0.1:8000/implicits", nil)
	c.Assert(err, IsNil)
	resp = httptest.NewRecorder()
	handler(resp, req)
	body2, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)

	c.Assert(string(body1), Equals, string(body2))
}
