package routing

import (
	"github.com/redneckbeard/gadget/controller"
	"github.com/redneckbeard/gadget/requests"
	. "launchpad.net/gocheck"
	"net/http"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type RouteSuite struct{}

func (s *RouteSuite) SetUpSuite(c *C) {
	controller.Register(&URLParamController{controller.New()})
	controller.Register(&TellMethodNameController{controller.New()})
}

var _ = Suite(&RouteSuite{})

type TellMethodNameController struct {
	*controller.DefaultController
}

func (c *TellMethodNameController) Index(r *requests.Request) (int, interface{}) {
	return 200, "index"
}

func (c *TellMethodNameController) Show(r *requests.Request) (int, interface{}) {
	return 200, "show"
}

func (c *TellMethodNameController) Create(r *requests.Request) (int, interface{}) {
	return 200, "create"
}

func (c *TellMethodNameController) Update(r *requests.Request) (int, interface{}) {
	return 200, "update"
}

func (c *TellMethodNameController) Destroy(r *requests.Request) (int, interface{}) {
	return 200, "destroy"
}

func (c *TellMethodNameController) Arbitrary(r *requests.Request) (int, interface{}) {
	return 200, "arbitrary"
}

//Route.Respond calls a controller's Index method on a GET request that matches the indexPattern
func (s *RouteSuite) TestRouterespondCallsControllersIndexMethodOnGetRequestThatMatchesIndexpattern(c *C) {
	r := newRoute("tellmethodname")
	r.buildPatterns("")
	req, _ := http.NewRequest("GET", "http://127.0.0.1:8000/tellmethodname", nil)
	status, body, action := r.Respond(requests.New(req))
	c.Assert(status, Equals, 200)
	c.Assert(body.(string), Equals, "index")
	c.Assert(action, Equals, "index")
}

//Route.Respond calls a controller's Show method on a GET request that matches the objectPattern
func (s *RouteSuite) TestRouterespondCallsControllersShowMethodOnGetRequestThatMatchesObjectpattern(c *C) {
	r := newRoute("tellmethodname")
	r.buildPatterns("")
	req, _ := http.NewRequest("GET", "http://127.0.0.1:8000/tellmethodname/1", nil)
	status, body, action := r.Respond(requests.New(req))
	c.Assert(status, Equals, 200)
	c.Assert(body.(string), Equals, "show")
	c.Assert(action, Equals, "show")
}

//Route.Respond 404s on a POST request that matches its controller's objectPattern
func (s *RouteSuite) TestRouterespond404SOnPostRequestThatMatchesItsControllersIndexpattern(c *C) {
	r := newRoute("tellmethodname")
	r.buildPatterns("")
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8000/tellmethodname/1", nil)
	status, body, action := r.Respond(requests.New(req))
	c.Assert(status, Equals, 404)
	c.Assert(body.(string), Equals, "")
	c.Assert(action, Equals, "")
}

//Route.Respond calls a controller's Create method on a POST request that matches the indexPattern
func (s *RouteSuite) TestRouterespondCallsControllersCreateMethodOnPostRequestThatMatchesObjectpattern(c *C) {
	r := newRoute("tellmethodname")
	r.buildPatterns("")
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8000/tellmethodname", nil)
	status, body, action := r.Respond(requests.New(req))
	c.Assert(status, Equals, 200)
	c.Assert(body.(string), Equals, "create")
	c.Assert(action, Equals, "create")
}

//Route.Respond 404s on a PUT request that matches its controller's indexPattern
func (s *RouteSuite) TestRouterespond404SOnPutRequestThatMatchesItsControllersIndexpattern(c *C) {
	r := newRoute("tellmethodname")
	r.buildPatterns("")
	req, _ := http.NewRequest("PUT", "http://127.0.0.1:8000/tellmethodname", nil)
	status, body, action := r.Respond(requests.New(req))
	c.Assert(status, Equals, 404)
	c.Assert(body.(string), Equals, "")
	c.Assert(action, Equals, "")
}

//Route.Respond calls a controller's Update method on a PUT request that matches the objectPattern
func (s *RouteSuite) TestRouterespondCallsControllersUpdateMethodOnPutRequestThatMatchesObjectpattern(c *C) {
	r := newRoute("tellmethodname")
	r.buildPatterns("")
	req, _ := http.NewRequest("PUT", "http://127.0.0.1:8000/tellmethodname/1", nil)
	status, body, action := r.Respond(requests.New(req))
	c.Assert(status, Equals, 200)
	c.Assert(body.(string), Equals, "update")
	c.Assert(action, Equals, "update")
}

//Route.Respond 404s on a DELETE request that matches its controller's indexPattern
func (s *RouteSuite) TestRouterespond404SOnDeleteRequestThatMatchesItsControllersIndexpattern(c *C) {
	r := newRoute("tellmethodname")
	r.buildPatterns("")
	req, _ := http.NewRequest("DELETE", "http://127.0.0.1:8000/tellmethodname", nil)
	status, body, action := r.Respond(requests.New(req))
	c.Assert(status, Equals, 404)
	c.Assert(body.(string), Equals, "")
	c.Assert(action, Equals, "")
}

//Route.Respond calls a controller's Destroy method on a DELETE request that matches the objectPattern
func (s *RouteSuite) TestRouterespondCallsControllersDestroyMethodOnDeleteRequestThatMatchesObjectpattern(c *C) {
	r := newRoute("tellmethodname")
	r.buildPatterns("")
	req, _ := http.NewRequest("DELETE", "http://127.0.0.1:8000/tellmethodname/1", nil)
	status, body, action := r.Respond(requests.New(req))
	c.Assert(status, Equals, 200)
	c.Assert(body.(string), Equals, "destroy")
	c.Assert(action, Equals, "destroy")
}

//Route.Respond should call an arbitrary exported method when a request path matches its name 
func (s *RouteSuite) TestRouterespondCallsControllersArbitraryMethodOnDeleteRequestThatMatchesObjectpattern(c *C) {
	r := newRoute("tellmethodname")
	r.buildPatterns("")
	req, _ := http.NewRequest("GET", "http://127.0.0.1:8000/tellmethodname/arbitrary", nil)
	status, body, action := r.Respond(requests.New(req))
	c.Assert(status, Equals, 200)
	c.Assert(body.(string), Equals, "arbitrary")
	c.Assert(action, Equals, "arbitrary")
}

type URLParamController struct {
	*controller.DefaultController
}

func (c *URLParamController) Index(r *requests.Request) (int, interface{}) {
	return 200, r.UrlParams["urlparam_id"]
}

func (c *URLParamController) Show(r *requests.Request) (int, interface{}) {
	return 200, r.UrlParams["urlparam_id"]
}

func (c *URLParamController) Create(r *requests.Request) (int, interface{}) {
	return 200, r.UrlParams["urlparam_id"]
}

func (c *URLParamController) Update(r *requests.Request) (int, interface{}) {
	return 200, r.UrlParams["urlparam_id"]
}

func (c *URLParamController) Destroy(r *requests.Request) (int, interface{}) {
	return 200, r.UrlParams["urlparam_id"]
}

//Route.Respond should pass UrlParams to objectPattern controller methods
func (s *RouteSuite) TestRouterespondShouldPassUrlparamsToObjectpatternControllerMethods(c *C) {
	r := newRoute("urlparam")
	r.buildPatterns("")
	objVerbs := []string{"GET", "PUT", "DELETE"}
	for _, verb := range objVerbs {
		req, _ := http.NewRequest(verb, "http://127.0.0.1:8000/urlparam/42", nil)
		status, body, _ := r.Respond(requests.New(req))
		c.Assert(status, Equals, 200)
		c.Assert(body.(string), Equals, "42")
	}
}

//Route.Respond should not pass UrlParams to indexPattern controller methods
func (s *RouteSuite) TestRouterespondShouldNotPassUrlparamsToIndexpatternControllerMethods(c *C) {
	r := newRoute("urlparam")
	r.buildPatterns("")
	idxVerbs := []string{"GET", "POST"}
	for _, verb := range idxVerbs {
		req, _ := http.NewRequest(verb, "http://127.0.0.1:8000/urlparam", nil)
		status, body, _ := r.Respond(requests.New(req))
		c.Assert(status, Equals, 200)
		c.Assert(body.(string), Equals, "")
	}
}
