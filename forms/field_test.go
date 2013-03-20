package forms

import (
	. "launchpad.net/gocheck"
	"errors"
	"fmt"
	"time"
)

type FieldSuite struct{}

var _ = Suite(&FieldSuite{})

//A StringField should return an error when cleaning an int
func (s *FieldSuite) TestStringfieldShouldReturnErrorWhenCleaningInt(c *C) {
	field := &StringField{BaseField: NewBaseField()}
	field.Set(30)
	field.Clean()
	c.Assert(field.Error(), DeepEquals, errors.New(field.GetMessage("type")))
}

//A StringField should set a string as a value when cleaning a string
func (s *FieldSuite) TestStringfieldShouldSetStringAsValueWhenCleaningString(c *C) {
	field := &StringField{BaseField: NewBaseField()}
	field.Set("foo")
	field.Clean()
	c.Assert(field.Error(), IsNil)
	c.Assert(field.Value, Equals, "foo")
}

//An IntField should return an error when cleaning a string
func (s *FieldSuite) TestIntfieldShouldReturnErrorWhenCleaningString(c *C) {
	field := &IntField{BaseField: NewBaseField()}
	field.Set("foo")
	field.Clean()
	c.Assert(field.Error(), DeepEquals, errors.New(field.GetMessage("type")))
}

//An IntField should set an int as a value when cleaning a string that can be coerced to an int
func (s *FieldSuite) TestIntfieldShouldSetIntAsValueWhenCleaningStringThatCanBeCoercedToInt(c *C) {
	field := &IntField{BaseField: NewBaseField()}
	field.Set("39")
	field.Clean()
	c.Assert(field.Error(), IsNil)
	c.Assert(field.Value, Equals, 39)
}

//An IntField should return an error when cleaning a float
func (s *FieldSuite) TestIntfieldShouldReturnErrorWhenCleaningFloat(c *C) {
	field := &IntField{BaseField: NewBaseField()}
	field.Set(39.5)
	field.Clean()
	c.Assert(field.Error(), DeepEquals, errors.New(field.GetMessage("type")))
}

//An IntField should set an int as a value when cleaning an int
func (s *FieldSuite) TestIntfieldShouldSetIntAsValueWhenCleaningInt(c *C) {
	field := &IntField{BaseField: NewBaseField()}
	field.Set(39)
	field.Clean()
	c.Assert(field.Error(), IsNil)
	c.Assert(field.Value, Equals, 39)
}

//A Float64Field should return an error when cleaning a string
func (s *FieldSuite) TestFloat64FieldShouldReturnErrorWhenCleaningString(c *C) {
	field := &Float64Field{BaseField: NewBaseField()}
	field.Set("foo")
	field.Clean()
	c.Assert(field.Error(), DeepEquals, errors.New(field.GetMessage("type")))
}

//A Float64Field should set a float as a value when cleaning a string that can be coerced to float64
func (s *FieldSuite) TestFloat64FieldShouldSetFloatAsValueWhenCleaningStringThatCanBeCoercedToFloat64(c *C) {
	field := &Float64Field{BaseField: NewBaseField()}
	field.Set("78.2")
	field.Clean()
	c.Assert(field.Error(), IsNil)
	c.Assert(field.Value, Equals, 78.2)
}

//A Float64Field should set a float as a value when cleaning an int
func (s *FieldSuite) TestFloat64FieldShouldSetFloatAsValueWhenCleaningInt(c *C) {
	field := &Float64Field{BaseField: NewBaseField()}
	field.Set(78)
	field.Clean()
	c.Assert(field.Error(), IsNil)
	c.Assert(field.Value, Equals, 78.0)
}

//A Float64Field should set a float as a value when cleaning a float64
func (s *FieldSuite) TestFloat64FieldShouldSetFloatAsValueWhenCleaningFloat64(c *C) {
	field := &Float64Field{BaseField: NewBaseField()}
	field.Set(82.7)
	field.Clean()
	c.Assert(field.Error(), IsNil)
	c.Assert(field.Value, Equals, 82.7)
}

//A BoolField should set false as its value when cleaning false
func (s *FieldSuite) TestBoolfieldShouldSetFalseAsItsValueWhenCleaningFalse(c *C) {
	field := &BoolField{BaseField: NewBaseField()}
	field.Set(false)
	field.Clean()
	c.Assert(field.Error(), IsNil)
	c.Assert(field.Value, Equals, false)
}

//A BoolField should set false as its value when cleaning 0
func (s *FieldSuite) TestBoolfieldShouldSetFalseAsItsValueWhenCleaning0(c *C) {
	field := &BoolField{BaseField: NewBaseField()}
	field.Set(0)
	field.Clean()
	c.Assert(field.Error(), IsNil)
	c.Assert(field.Value, Equals, false)
}

//A BoolField should set false as its value when cleaning an empty string
func (s *FieldSuite) TestBoolfieldShouldSetFalseAsItsValueWhenCleaningEmptyString(c *C) {
	field := &BoolField{BaseField: NewBaseField()}
	field.Set("")
	field.Clean()
	c.Assert(field.Error(), IsNil)
	c.Assert(field.Value, Equals, false)
}

//A BoolField should set true as its value when cleaning true
func (s *FieldSuite) TestBoolfieldShouldSetTrueAsItsValueWhenCleaningTrue(c *C) {
	field := &BoolField{BaseField: NewBaseField()}
	field.Set(true)
	field.Clean()
	c.Assert(field.Error(), IsNil)
	c.Assert(field.Value, Equals, true)
}

//A BoolField should set true as its value when cleaning 1
func (s *FieldSuite) TestBoolfieldShouldSetTrueAsItsValueWhenCleaning1(c *C) {
	field := &BoolField{BaseField: NewBaseField()}
	field.Set(1)
	field.Clean()
	c.Assert(field.Error(), IsNil)
	c.Assert(field.Value, Equals, true)
}

//A BoolField should set true as its value when cleaning "1"
func (s *FieldSuite) TestBoolfieldShouldSetTrueAsItsValueWhenCleaningString1(c *C) {
	field := &BoolField{BaseField: NewBaseField()}
	field.Set("1")
	field.Clean()
	c.Assert(field.Error(), IsNil)
	c.Assert(field.Value, Equals, true)
}

//A TimeField should return an error when cleaning a negative integer
func (s *FieldSuite) TestTimefieldShouldReturnErrorWhenCleaningNegativeInteger(c *C) {
	field := &TimeField{BaseField: NewBaseField()}
	field.DefaultMessages()
	field.Set(-100101)
	field.Clean()
	c.Assert(field.Error(), ErrorMatches, fmt.Sprintf(field.GetMessage("type"), field.Format))
}

//A TimeField should set a time.Time as the seconds from the epoch when cleaning a positive integer
func (s *FieldSuite) TestTimefieldShouldSetTimetimeAsSecondsFromEpochWhenCleaningPositiveInteger(c *C) {
	t := time.Now()
	t = t.Add(time.Duration(-(t.Nanosecond()))).UTC()
	field := &TimeField{BaseField: NewBaseField()}
	field.Set(t.Unix())
	field.Clean()
	c.Assert(field.Error(), IsNil)
	c.Assert(field.Value, Equals, t)
}

//A TimeField should return an error when cleaning a string that does not match the field's format string
func (s *FieldSuite) TestTimefieldShouldReturnErrorWhenCleaningStringThatDoesNotMatchFieldsFormatString(c *C) {
	t := time.Now()
	field := &TimeField{BaseField: NewBaseField()}
	field.DefaultMessages()
	field.Set(t.Format(time.RFC822))
	field.Clean()
	c.Assert(field.Error(), ErrorMatches, fmt.Sprintf(field.GetMessage("type"), field.Format))
}

//A TimeField should set a time.Time when cleaning a string that matches the field's format string
func (s *FieldSuite) TestTimefieldShouldSetTimetimeWhenCleaningStringThatMatchesFieldsFormatString(c *C) {
	t := time.Now()
	t = t.Add(time.Duration(-(t.Nanosecond()))).UTC()
	field := &TimeField{BaseField: NewBaseField()}
	field.Format = time.ANSIC
	field.Set(t.Format(time.ANSIC))
	field.Clean()
	c.Assert(field.Error(), IsNil)
	c.Assert(field.Value, Equals, t)
}
