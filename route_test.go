package gadget

import (
	. "launchpad.net/gocheck"
	"net/http"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type RouteSuite struct{}

type routeApp struct { *App }

var rta *routeApp

func (s *RouteSuite) SetUpTest(c *C) {
	rta = &routeApp{ &App{} }
	rta.Register(&URLParamController{})
	rta.Register(&TellMethodNameController{})
}

func (s *RouteSuite) TearDownTest(c *C) {
	rta.Controllers = make(map[string]Controller)
}

var _ = Suite(&RouteSuite{})

type TellMethodNameController struct {
	*DefaultController
}

func (c *TellMethodNameController) Index(r *Request) (int, interface{}) {
	return 200, "index"
}

func (c *TellMethodNameController) Show(r *Request) (int, interface{}) {
	return 200, "show"
}

func (c *TellMethodNameController) Create(r *Request) (int, interface{}) {
	return 200, "create"
}

func (c *TellMethodNameController) Update(r *Request) (int, interface{}) {
	return 200, "update"
}

func (c *TellMethodNameController) Destroy(r *Request) (int, interface{}) {
	return 200, "destroy"
}

func (c *TellMethodNameController) Arbitrary(r *Request) (int, interface{}) {
	return 200, "arbitrary"
}

//Route.Respond calls a controller's Index method on a GET request that matches the indexPattern
func (s *RouteSuite) TestRouterespondCallsControllersIndexMethodOnGetRequestThatMatchesIndexpattern(c *C) {
	r := rta.newRoute("tell-method-names", nil)
	r.buildPatterns("")
	req, _ := http.NewRequest("GET", "http://127.0.0.1:8000/tell-method-names", nil)
	status, body, action := r.Respond(newRequest(req))
	c.Assert(status, Equals, 200)
	c.Assert(body.(string), Equals, "index")
	c.Assert(action, Equals, "index")
}

//Route.Respond calls a controller's Show method on a GET request that matches the objectPattern
func (s *RouteSuite) TestRouterespondCallsControllersShowMethodOnGetRequestThatMatchesObjectpattern(c *C) {
	r := rta.newRoute("tell-method-names", nil)
	r.buildPatterns("")
	req, _ := http.NewRequest("GET", "http://127.0.0.1:8000/tell-method-names/1", nil)
	status, body, action := r.Respond(newRequest(req))
	c.Assert(status, Equals, 200)
	c.Assert(body.(string), Equals, "show")
	c.Assert(action, Equals, "show")
}

//Route.Respond 404s on a POST request that matches its controller's objectPattern
func (s *RouteSuite) TestRouterespond404SOnPostRequestThatMatchesItsControllersIndexpattern(c *C) {
	r := rta.newRoute("tell-method-names", nil)
	r.buildPatterns("")
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8000/tell-method-names/1", nil)
	status, body, action := r.Respond(newRequest(req))
	c.Assert(status, Equals, 404)
	c.Assert(body.(string), Equals, "")
	c.Assert(action, Equals, "")
}

//Route.Respond calls a controller's Create method on a POST request that matches the indexPattern
func (s *RouteSuite) TestRouterespondCallsControllersCreateMethodOnPostRequestThatMatchesObjectpattern(c *C) {
	r := rta.newRoute("tell-method-names", nil)
	r.buildPatterns("")
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8000/tell-method-names", nil)
	status, body, action := r.Respond(newRequest(req))
	c.Assert(status, Equals, 200)
	c.Assert(body.(string), Equals, "create")
	c.Assert(action, Equals, "create")
}

//Route.Respond 404s on a PUT request that matches its controller's indexPattern
func (s *RouteSuite) TestRouterespond404SOnPutRequestThatMatchesItsControllersIndexpattern(c *C) {
	r := rta.newRoute("tell-method-names", nil)
	r.buildPatterns("")
	req, _ := http.NewRequest("PUT", "http://127.0.0.1:8000/tell-method-names", nil)
	status, body, action := r.Respond(newRequest(req))
	c.Assert(status, Equals, 404)
	c.Assert(body.(string), Equals, "")
	c.Assert(action, Equals, "")
}

//Route.Respond calls a controller's Update method on a PUT request that matches the objectPattern
func (s *RouteSuite) TestRouterespondCallsControllersUpdateMethodOnPutRequestThatMatchesObjectpattern(c *C) {
	r := rta.newRoute("tell-method-names", nil)
	r.buildPatterns("")
	req, _ := http.NewRequest("PUT", "http://127.0.0.1:8000/tell-method-names/1", nil)
	status, body, action := r.Respond(newRequest(req))
	c.Assert(status, Equals, 200)
	c.Assert(body.(string), Equals, "update")
	c.Assert(action, Equals, "update")
}

//Route.Respond calls a controller's Update method on a PATCH request that matches the objectPattern
func (s *RouteSuite) TestRouterespondCallsControllersUpdateMethodOnPatchRequestThatMatchesObjectpattern(c *C) {
	r := rta.newRoute("tell-method-names", nil)
	r.buildPatterns("")
	req, _ := http.NewRequest("PATCH", "http://127.0.0.1:8000/tell-method-names/1", nil)
	status, body, action := r.Respond(newRequest(req))
	c.Assert(status, Equals, 200)
	c.Assert(body.(string), Equals, "update")
	c.Assert(action, Equals, "update")
}

//Route.Respond 404s on a DELETE request that matches its controller's indexPattern
func (s *RouteSuite) TestRouterespond404SOnDeleteRequestThatMatchesItsControllersIndexpattern(c *C) {
	r := rta.newRoute("tell-method-names", nil)
	r.buildPatterns("")
	req, _ := http.NewRequest("DELETE", "http://127.0.0.1:8000/tell-method-names", nil)
	status, body, action := r.Respond(newRequest(req))
	c.Assert(status, Equals, 404)
	c.Assert(body.(string), Equals, "")
	c.Assert(action, Equals, "")
}

//Route.Respond calls a controller's Destroy method on a DELETE request that matches the objectPattern
func (s *RouteSuite) TestRouterespondCallsControllersDestroyMethodOnDeleteRequestThatMatchesObjectpattern(c *C) {
	r := rta.newRoute("tell-method-names", nil)
	r.buildPatterns("")
	req, _ := http.NewRequest("DELETE", "http://127.0.0.1:8000/tell-method-names/1", nil)
	status, body, action := r.Respond(newRequest(req))
	c.Assert(status, Equals, 200)
	c.Assert(body.(string), Equals, "destroy")
	c.Assert(action, Equals, "destroy")
}

//Route.Respond should call an arbitrary exported method when a request path matches its name
func (s *RouteSuite) TestRouterespondCallsControllersArbitraryMethodOnDeleteRequestThatMatchesObjectpattern(c *C) {
	r := rta.newRoute("tell-method-names", nil)
	r.buildPatterns("")
	req, _ := http.NewRequest("GET", "http://127.0.0.1:8000/tell-method-names/arbitrary", nil)
	status, body, action := r.Respond(newRequest(req))
	c.Assert(status, Equals, 200)
	c.Assert(body.(string), Equals, "arbitrary")
	c.Assert(action, Equals, "arbitrary")
}

type URLParamController struct {
	*DefaultController
}

func (c *URLParamController) Index(r *Request) (int, interface{}) {
	return 200, r.UrlParams["urlparam_id"]
}

func (c *URLParamController) Show(r *Request) (int, interface{}) {
	return 200, r.UrlParams["urlparam_id"]
}

func (c *URLParamController) Create(r *Request) (int, interface{}) {
	return 200, r.UrlParams["urlparam_id"]
}

func (c *URLParamController) Update(r *Request) (int, interface{}) {
	return 200, r.UrlParams["urlparam_id"]
}

func (c *URLParamController) Destroy(r *Request) (int, interface{}) {
	return 200, r.UrlParams["urlparam_id"]
}

//Route.Respond should pass UrlParams to objectPattern controller methods
func (s *RouteSuite) TestRouterespondShouldPassUrlparamsToObjectpatternControllerMethods(c *C) {
	r := rta.newRoute("urlparams", nil)
	r.buildPatterns("")
	objVerbs := []string{"GET", "PUT", "DELETE"}
	for _, verb := range objVerbs {
		req, _ := http.NewRequest(verb, "http://127.0.0.1:8000/urlparams/42", nil)
		status, body, _ := r.Respond(newRequest(req))
		c.Assert(status, Equals, 200)
		c.Assert(body.(string), Equals, "42")
	}
}

//Route.Respond should not pass UrlParams to indexPattern controller methods
func (s *RouteSuite) TestRouterespondShouldNotPassUrlparamsToIndexpatternControllerMethods(c *C) {
	r := rta.newRoute("urlparams", nil)
	r.buildPatterns("")
	idxVerbs := []string{"GET", "POST"}
	for _, verb := range idxVerbs {
		req, _ := http.NewRequest(verb, "http://127.0.0.1:8000/urlparams", nil)
		status, body, _ := r.Respond(newRequest(req))
		c.Assert(status, Equals, 200)
		c.Assert(body.(string), Equals, "")
	}
}

func AclFilter(r *Request) (status int, body interface{}) {
	if r.UrlParams["urlparam_id"] == "10" {
		status, body = 403, "VERBOTEN"
	}
	return
}

func (s *RouteSuite) TestFilterReturnValueUsed(c *C) {
	ctrl, _ := rta.GetController("urlparams")
	ctrl.Filter([]string{"update"}, AclFilter)
	r := rta.newRoute("urlparams", nil)
	r.buildPatterns("")
	req, _ := http.NewRequest("PUT", "http://127.0.0.1:8000/urlparams/10", nil)
	status, body, _ := r.Respond(newRequest(req))
	c.Assert(status, Equals, 403)
	c.Assert(body.(string), Equals, "VERBOTEN")
}

func (s *RouteSuite) TestFilterOnlyAppliedToSpecifiedActions(c *C) {
	ctrl, _ := rta.GetController("urlparams")
	ctrl.Filter([]string{"update"}, AclFilter)
	r := rta.newRoute("urlparams", nil)
	r.buildPatterns("")
	req, _ := http.NewRequest("GET", "http://127.0.0.1:8000/urlparams/10", nil)
	status, body, _ := r.Respond(newRequest(req))
	c.Assert(status, Equals, 200)
	c.Assert(body.(string), Equals, "10")
}

func (s *RouteSuite) TestRespondContinuesAfterNilFilterReturn(c *C) {
	ctrl, _ := rta.GetController("urlparams")
	ctrl.Filter([]string{"update"}, AclFilter)
	r := rta.newRoute("urlparams", nil)
	r.buildPatterns("")
	req, _ := http.NewRequest("PUT", "http://127.0.0.1:8000/urlparams/11", nil)
	status, body, _ := r.Respond(newRequest(req))
	c.Assert(status, Equals, 200)
	c.Assert(body.(string), Equals, "11")
}
