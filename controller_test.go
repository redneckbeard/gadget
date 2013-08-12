package gadget

import (
	. "launchpad.net/gocheck"
	"github.com/redneckbeard/gadget/strutil"
)

type ControllerSuite struct{}

func (s *ControllerSuite) SetUpTest(c *C) {
	Register(&TestController{})
}

func (s *ControllerSuite) TearDownTest(c *C) {
	clear()
}

var _ = Suite(&ControllerSuite{})

type TestController struct {
	*DefaultController
}

type BrokenName struct {
	*DefaultController
}

//hyphenate should convert "PascalCase" into "pascal-case"
func (s *ControllerSuite) TestHyphenateShouldConvertPascalcaseIntoPascalcase(c *C) {
	c.Assert(strutil.Hyphenate("PascalCase"), Equals, "pascal-case")
}

//NameOf(&TestController) should return "test"
func (s *ControllerSuite) TestNameoftestcontrollerShouldReturnTest(c *C) {
	t := &TestController{}
	c.Assert(nameOf(t), Equals, "test")
}

//NameOf(&BrokenName) should panic
func (s *ControllerSuite) TestNameofbrokennameShouldPanic(c *C) {
	b := &BrokenName{}
	c.Assert(func() { nameOf(b) }, PanicMatches, `Controller names must adhere to the convention of '<name>Controller'`)
}

//Get("test") should return a *TestController if one is registered
func (s *ControllerSuite) TestGettestShouldReturnTestcontrollerOneIsRegistered(c *C) {
	ctlr, _ := getController("tests")
	_, ok := ctlr.(*TestController)
	c.Assert(ok, Equals, true)
}

//Get("other") should return an error if no such controller is registered
func (s *ControllerSuite) TestGetotherShouldReturnErrorNoSuchControllerIsRegistered(c *C) {
	_, err := getController("other")
	c.Assert(err, ErrorMatches, "No controller with label 'other' found")
}

func F(r *Request) (int, interface{}) {
	return 200, "OK"
}

//controller.Filter("missing", Filter) should panic if the controller has no method "missing"
func (s *ControllerSuite) TestControllerfiltermissingFilterShouldPanicControllerHasNoMethodMissing(c *C) {
	ctlr, _ := getController("tests")
	c.Assert(func() { ctlr.Filter([]string{"missing"}, F) }, PanicMatches, "Unable to add filter for 'missing' -- no such action")
}
