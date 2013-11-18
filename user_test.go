package gadget

import (
	"io/ioutil"
	. "launchpad.net/gocheck"
	"net/http"
	"net/http/httptest"
)

type UserSuite struct{}

type userApp struct {
	*App
}

var u *userApp

func (s *UserSuite) SetUpTest(c *C) {
	u = &userApp{ &App{} }
	u.Register(&AuthStatusController{})
	u.Routes(u.Resource("auth-status"))
}

func (s *UserSuite) TearDownTest(c *C) {
	clearUserIdentifier()
}

var _ = Suite(&UserSuite{})

type AuthStatusController struct {
	*DefaultController
}

func (c *AuthStatusController) Plural() string { return "auth-status" }

func (c *AuthStatusController) Index(r *Request) (int, interface{}) {
	return 200, r.User.Authenticated()
}

type AuthedUser struct{}

func (u *AuthedUser) Authenticated() bool { return true }

func FakeAuth(r *Request) User {
	if authed, ok := r.Params["authed"]; ok {
		if authed.(string) == "yes" {
			return &AuthedUser{}
		}
	}
	return &AnonymousUser{}
}

//The User attached to the Request should not be authenticated if no UserIdentifier has been registered
func (s *UserSuite) TestUserAttachedToRequestAuthenticatedNoUseridentifierHasBeenRegistered(c *C) {
	handler := u.Handler()

	req, _ := http.NewRequest("GET", "http://127.0.0.1:8000/auth-status", nil)
	resp := httptest.NewRecorder()
	handler(resp, req)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "false")
}

//The User attached to the Request should not be authenticated if the registered UserIdentifier returns an anonymous user
func (s *UserSuite) TestUserAttachedToRequestAuthenticatedRegisteredUseridentifierReturnsAnonymousUser(c *C) {
	handler := u.Handler()

	IdentifyUsersWith(FakeAuth)
	req, _ := http.NewRequest("GET", "http://127.0.0.1:8000/auth-status", nil)
	resp := httptest.NewRecorder()
	handler(resp, req)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "false")
}

//The User attached to the Request should be authenticated if the registered UserIdentifier returns an authenticated user
func (s *UserSuite) TestUserAttachedToRequestAuthenticatedRegisteredUseridentifierReturnsAuthenticatedUser(c *C) {
	handler := u.Handler()

	IdentifyUsersWith(FakeAuth)
	req, _ := http.NewRequest("GET", "http://127.0.0.1:8000/auth-status?authed=yes", nil)
	resp := httptest.NewRecorder()
	handler(resp, req)
	body, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, "true")
}
