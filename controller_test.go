package gadget

import (
	. "launchpad.net/gocheck"
)

type ControllerSuite struct{}

func (s *ControllerSuite) SetUpTest(c *C) {
	Register(&TestController{New()})
}

func (s *ControllerSuite) TearDownTest(c *C) {
	Clear()
}

var _ = Suite(&ControllerSuite{})

type TestController struct {
	*DefaultController
}

type BrokenName struct {
	*DefaultController
}

//NameOf(&TestController) should return "test"
func (s *ControllerSuite) TestNameoftestcontrollerShouldReturnTest(c *C) {
	t := &TestController{New()}
	c.Assert(NameOf(t), Equals, "test")
}

//NameOf(&BrokenName) should panic
func (s *ControllerSuite) TestNameofbrokennameShouldPanic(c *C) {
	b := &BrokenName{New()}
	c.Assert(func() { NameOf(b) }, PanicMatches, `Controller names must adhere to the convention of '<name>Controller'`)
}

//Get("test") should return a *TestController if one is registered
func (s *ControllerSuite) TestGettestShouldReturnTestcontrollerOneIsRegistered(c *C) {
	ctlr, _ := Get("test")
	_, ok := ctlr.(*TestController)
	c.Assert(ok, Equals, true)
}

//Get("other") should return an error if no such controller is registered
func (s *ControllerSuite) TestGetotherShouldReturnErrorNoSuchControllerIsRegistered(c *C) {
	_, err := Get("other")
	c.Assert(err, ErrorMatches, "No controller with label 'other' found")
}

func F(r *Request) (int, interface{}) {
		return 200, "OK"
	}


//controller.Filter("missing", Filter) should panic if the controller has no method "missing"
func (s *ControllerSuite) TestControllerfiltermissingFilterShouldPanicControllerHasNoMethodMissing(c *C) {
	ctlr, _ := Get("test")
	c.Assert(func() { ctlr.Filter([]string{"missing"}, F) }, PanicMatches, "Unable to add filter for 'missing' -- no such action")
}
