package gadget

import (
	"io/ioutil"
	. "launchpad.net/gocheck"
	"net/http"
	"net/http/httptest"
	"time"
)

type ResponseSuite struct{}

type responseApp struct {
	*App
}

var ra *responseApp

var _ = Suite(&ResponseSuite{})

func (s *ResponseSuite) SetUpSuite(c *C) {
	ra = &responseApp{&App{}}
	ra.Register(&ResponseController{})
	ra.Register(&ImplicitController{})
	ra.Accept("application/json").Via(JsonBroker)
	ra.Accept("text/html").Via(HtmlBroker)
	ra.Routes(ra.Resource("responses"), ra.Resource("implicits"))
}
func (s *ResponseSuite) TearDownSuite(c *C) {
	ra.Controllers = make(map[string]Controller)
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

func (c *ResponseController) CookieAndRedirect(*Request) (int, interface{}) {
	response := NewResponse("/responses")
	response.AddCookie(cookie)
	return 302, response
}

func (c *ResponseController) RedirectWithString(*Request) (int, interface{}) {
	return 301, "/somewhere"
}

type ImplicitController struct {
	*DefaultController
}

func (c *ImplicitController) Index(*Request) (int, interface{}) {
	body := struct{ Foo, Bar string }{"baz", "quux"}
	return 200, body
}

func HtmlBroker(r *Request, status int, body interface{}, data *RouteData) (int, string) {
	return 200, ""
}

//Headers set on a Response in a controller method are correctly transferred to the http.Response
func (s *ResponseSuite) TestHeadersSetOnGadgetresponseControllerMethodAreCorrectlyTransferredToHttpresponse(c *C) {
	handler := ra.Handler()

	req, err := http.NewRequest("GET", "http://127.0.0.1:8000/responses", nil)
	c.Assert(err, IsNil)

	resp := httptest.NewRecorder()
	handler(resp, req)
	c.Assert(resp.Header().Get("X-Framework"), Equals, "Gadget")
}

//Cookies added to a Response in a controller method are correctly transferred to the http.Response
func (s *ResponseSuite) TestCookiesAddedToGadgetresponseControllerMethodAreCorrectlyTransferredToHttpresponse(c *C) {
	handler := ra.Handler()

	req, err := http.NewRequest("GET", "http://127.0.0.1:8000/responses/1", nil)
	c.Assert(err, IsNil)

	resp := httptest.NewRecorder()
	handler(resp, req)
	c.Assert(resp.Code, Equals, 200)
	c.Assert(resp.Header().Get("Set-Cookie"), Equals, cookie.String())
}

//The body of the http.Response is identical between a controller method that returns an anonymous struct and one that returns a Response with its Body set to that struct
func (s *ResponseSuite) TestBodyHttpresponseIsIdenticalBetweenControllerMethodThatReturnsAnonymousStructAndOneThatReturnsGadgetresponseItsBodySetToThatStruct(c *C) {
	handler := ra.Handler()

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

// The Content-Type of the outgoing response should be the same as the incoming request provided our app supports that mime type.
func (s *ResponseSuite) TestContentTypeOfCustomHttpResponseMatchesRequestContentType(c *C) {
	handler := ra.Handler()

	req, err := http.NewRequest("GET", "http://127.0.0.1:8000/responses", nil)
	req.Header.Add("Content-Type", "text/html")
	c.Assert(err, IsNil)
	resp := httptest.NewRecorder()
	handler(resp, req)

	c.Assert(resp.Header().Get("Content-Type"), Equals, req.Header.Get("Content-Type"))
}

// We should be able to set a cookie and redirect using the technique outlined in the CookieAndRedirect action.
func (s *ResponseSuite) TestSetCookieAndRedirect(c *C) {
	handler := ra.Handler()

	req, err := http.NewRequest("GET", "http://127.0.0.1:8000/responses/cookie-and-redirect", nil)
	c.Assert(err, IsNil)
	resp := httptest.NewRecorder()
	handler(resp, req)

	c.Assert(resp.Code, Equals, 302)
	c.Assert(resp.Header().Get("Set-Cookie"), Equals, cookie.String())
}

func (s *ResponseSuite) TestRedirectWithString(c *C) {
	handler := ra.Handler()

	req, err := http.NewRequest("GET", "http://127.0.0.1:8000/responses/redirect-with-string", nil)
	c.Assert(err, IsNil)
	resp := httptest.NewRecorder()
	handler(resp, req)

	c.Assert(resp.Code, Equals, 301)
	c.Assert(resp.Header().Get("Location"), Equals, "/somewhere")
}
