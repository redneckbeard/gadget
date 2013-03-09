package gadget

import (
	. "launchpad.net/gocheck"
)

type RegistrySuite struct{}

func (s *RegistrySuite) SetUpSuite(c *C) {
	Register(&FooController{New()})
	Register(&BarController{New()})
	Register(&BazController{New()})
}

var _ = Suite(&RegistrySuite{})

type FooController struct {
	*DefaultController
}

type BarController struct {
	*DefaultController
}

type BazController struct {
	*DefaultController
}

//SetIndex("foo") should return a Route with no subroutes and an indexPattern of ^$
func (s *RegistrySuite) TestRoutingindexfooShouldReturnRouteNoSubroutesAndIndexpattern(c *C) {
	r := SetIndex("foo")
	c.Assert(r.subroutes, HasLen, 0)
	c.Assert(r.indexPattern.String(), Equals, `^$`)
}

//Resource("foo") should have no subroutes 
func (s *RegistrySuite) TestRoutingresourcefooHasNoSubroutesAndIndexpatternFoo(c *C) {
	r := Resource("foo")
	c.Assert(r.subroutes, HasLen, 0)
}

//Resource("foo") should have an indexPattern of ^foo$
func (s *RegistrySuite) TestRoutingresourcefooHasIndexpatternFoo(c *C) {
	r := Resource("foo")
	c.Assert(r.indexPattern.String(), Equals, `^foo$`)
}

//Resource("foo") should have an objectPattern of ^foo(?:/(?P<foo_id>\d+))?$
func (s *RegistrySuite) TestRoutingresourcefooHasObjectpatternFoopfoo_Idd(c *C) {
	r := Resource("foo")
	c.Assert(r.objectPattern.String(), Equals, `^foo(?:/(?P<foo_id>\d+))?$`)
}

//Resource("foo", Resource("bar"), Resource("baz")) :
//* should return a Route with two subroutes
func (s *RegistrySuite) TestShouldReturnRouteTwoSubroutes(c *C) {
	r := Resource("foo", Resource("bar"), Resource("baz"))
	c.Assert(r.subroutes, HasLen, 2)
}

//* Calling Flatten on it should return a slice of 3 Routes
func (s *RegistrySuite) TestCallingFlattenOnShouldReturnSlice3RoutesNested(c *C) {
	r := Resource("foo", Resource("bar"), Resource("baz"))
	flattened := r.Flatten()
	c.Assert(flattened, HasLen, 3)
}

//* First route has indexPattern of ^foo$ and objectPattern of ^foo(?:/(?P<foo_id>\d+))?$
func (s *RegistrySuite) TestFirstRouteHasIndexpatternFooAndObjectpatternFoopfoo_Idd(c *C) {
	r := Resource("foo", Resource("bar"), Resource("baz"))
	first := r.Flatten()[0]
	c.Assert(first.indexPattern.String(), Equals, `^foo$`)
	c.Assert(first.objectPattern.String(), Equals, `^foo(?:/(?P<foo_id>\d+))?$`)
}

//* Second route has indexPattern of ^foo/(?P<foo_id>\d+)/bar$ and objectPattern of ^foo/(?<foo_id>\d+)/bar(?:/(?<bar_id>\d+))?$
func (s *RegistrySuite) TestSecondRouteHasIndexpatternFoopfoo_IddbarAndObjectpatternFoofoo_Iddbarbar_Idd(c *C) {
	r := Resource("foo", Resource("bar"), Resource("baz"))
	first := r.Flatten()[1]
	c.Assert(first.indexPattern.String(), Equals, `^foo/(?P<foo_id>\d+)/bar$`)
	c.Assert(first.objectPattern.String(), Equals, `^foo/(?P<foo_id>\d+)/bar(?:/(?P<bar_id>\d+))?$`)
}

//* Third route has indexPattern of ^foo/(?P<foo_id>\d+)/baz$ and objectPattern of ^foo/(?<foo_id>\d+)/baz(?:/(?<baz_id>\d+))?$
func (s *RegistrySuite) TestThirdRouteHasIndexpatternFoopfoo_IddbazAndObjectpatternFoofoo_Iddbazbaz_Idd(c *C) {
	r := Resource("foo", Resource("bar"), Resource("baz"))
	first := r.Flatten()[2]
	c.Assert(first.indexPattern.String(), Equals, `^foo/(?P<foo_id>\d+)/baz$`)
	c.Assert(first.objectPattern.String(), Equals, `^foo/(?P<foo_id>\d+)/baz(?:/(?P<baz_id>\d+))?$`)
}

//Resource("foo", Resource("bar", Resource("baz"))) :
//* should return a Route with a subroute with a subroute
func (s *RegistrySuite) TestShouldReturnRouteSubrouteSubroute(c *C) {
	r := Resource("foo", Resource("bar", Resource("baz")))
	c.Assert(r.subroutes, HasLen, 1)
	sub1 := r.subroutes[0]
	c.Assert(sub1.subroutes, HasLen, 1)
}

//* Calling Flatten on it should return a slice of 3 Routes
func (s *RegistrySuite) TestCallingFlattenOnShouldReturnSlice3Routes(c *C) {
	r := Resource("foo", Resource("bar", Resource("baz")))
	routes := r.Flatten()
	c.Assert(routes, HasLen, 3)
}

//* Third route has indexPattern of ^foo/(?P<foo_id>\d+)/bar/(?P<bar_id>\d+)/baz$ and objectPattern of ^foo/(?<foo_id>\d+)/bar/(?P<bar_id>\d+)/baz(?:/(?<baz_id>\d+))?$
func (s *RegistrySuite) TestThirdRouteHasIndexpatternFoopfoo_Iddbarpbar_IddbazAndObjectpatternFoofoo_Iddbarpbar_Iddbazbaz_Idd(c *C) {
	r := Resource("foo", Resource("bar", Resource("baz")))
	route := r.Flatten()[2]
	c.Assert(route.indexPattern.String(), Equals, `^foo/(?P<foo_id>\d+)/bar/(?P<bar_id>\d+)/baz$`)
	c.Assert(route.objectPattern.String(), Equals, `^foo/(?P<foo_id>\d+)/bar/(?P<bar_id>\d+)/baz(?:/(?P<baz_id>\d+))?$`)
}

//Prefixed("foo", Resource("bar"), Resource("baz")) :
//* should return a route with two subroutes
func (s *RegistrySuite) TestPrefixedShouldReturnRouteTwoSubroutes(c *C) {
	r := Prefixed("foo", Resource("bar"), Resource("baz"))
	c.Assert(r.subroutes, HasLen, 2)
}

//* calling Flatten should return a slice of 2 Routes
func (s *RegistrySuite) TestCallingFlattenShouldReturnSlice2Routes(c *C) {
	r := Prefixed("foo", Resource("bar"), Resource("baz"))
	subroutes := r.Flatten()
	c.Assert(subroutes, HasLen, 2)
}

//* First route has indexPattern of ^foo/bar$ and objectPattern of ^foo/bar(?:/(?P<bar_id>\d+))?$
func (s *RegistrySuite) TestFirstRouteHasIndexpatternFoobarAndObjectpatternFoobarpbar_Idd(c *C) {
	r := Prefixed("foo", Resource("bar"), Resource("baz"))
	subroute := r.Flatten()[0]
	c.Assert(subroute.indexPattern.String(), Equals, `^foo/bar$`)
	c.Assert(subroute.objectPattern.String(), Equals, `^foo/bar(?:/(?P<bar_id>\d+))?$`)
}

//* Second route has indexPattern of ^foo/baz$ and objectPattern of ^foo/baz(?:/(?P<baz_id>\d+))?$
func (s *RegistrySuite) TestSecondRouteHasIndexpatternFoobazAndObjectpatternFoobazpbaz_Idd(c *C) {
	r := Prefixed("foo", Resource("bar"), Resource("baz"))
	subroute := r.Flatten()[1]
	c.Assert(subroute.indexPattern.String(), Equals, `^foo/baz$`)
	c.Assert(subroute.objectPattern.String(), Equals, `^foo/baz(?:/(?P<baz_id>\d+))?$`)
}
