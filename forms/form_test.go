package forms

import (
	. "launchpad.net/gocheck"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type FormSuite struct{}

var _ = Suite(&FormSuite{})

type WidgetForm struct {
	*DefaultForm
	Foo, Bar  *StringField
	Baz, Quux *IntField
}

func (form *WidgetForm) CleanQuux(f FormField) error {
	field := f.(*IntField)
	i := field.Value
	if i < 35 {
		field.SetError("Quux values must be greater than 35.")
	}
	return nil
}

type BrokenWidget struct {
	Foo, Bar string
	Baz      int
}

type Widget struct {
	Foo, Bar  string
	Baz, Quux int
}

//Calling Ready() on an uninitialized Form should return false
func (s *FormSuite) TestCallingReadyOnUninitializedFormShouldReturnFalse(c *C) {
	f := &WidgetForm{}
	c.Assert(f.Ready(), Equals, false)
}

//Calling Ready() on an initialized Form should return true
func (s *FormSuite) TestCallingReadyOnInitializedFormShouldReturnFalse(c *C) {
	f := &WidgetForm{}
	Init(f)
	c.Assert(f.Ready(), Equals, true)
}

//An initialized Form should have a properly created Fields map
func (s *FormSuite) TestInitializedFormHasProperlyCreatedFieldsMap(c *C) {
	f := &WidgetForm{}
	Init(f)
	fields := f.FieldMap()
	c.Assert(fields["Foo"], FitsTypeOf, &StringField{BaseField: NewBaseField()})
	c.Assert(fields["Bar"], FitsTypeOf, &StringField{BaseField: NewBaseField()})
	c.Assert(fields["Bar"], FitsTypeOf, &StringField{BaseField: NewBaseField()})
	c.Assert(fields["Quux"], FitsTypeOf, &IntField{})
}

//Calling Populate on an uninitialized Form should return an error
func (s *FormSuite) TestCallingPopulateOnUninitializedFormShouldReturnError(c *C) {
	form := &WidgetForm{}
	err := Populate(form, map[string]interface{}{
		"Foo":  "baz",
		"Bar":  "quux",
		"Baz":  "37",
		"Quux": "52",
	})
	c.Assert(err, Not(IsNil))
}

//Calling Populate on a TestForm with valid data should set the values of those fields
func (s *FormSuite) TestCallingPopulateOnTestformValidDataShouldSetValuesThoseFields(c *C) {
	form := &WidgetForm{}
	Init(form)
	Populate(form, map[string]interface{}{
		"Foo":  "baz",
		"Bar":  "quux",
		"Baz":  "37",
		"Quux": "52",
	})
	c.Assert(form.Foo.Data.(string), DeepEquals, "baz")
	c.Assert(form.Bar.Data.(string), DeepEquals, "quux")
	c.Assert(form.Baz.Data.(string), DeepEquals, "37")
}

//Calling Populate on a TestForm with invalid data should set the values of those fields
func (s *FormSuite) TestCallingPopulateOnTestformInvalidDataShouldSetValuesThoseFields(c *C) {
	form := &WidgetForm{}
	Init(form)
	Populate(form, map[string]interface{}{
		"Foo":  "baz",
		"Bar":  "quux",
		"Baz":  "thirty-seven",
		"Quux": "52",
	})
	c.Assert(form.Foo.Data.(string), DeepEquals, "baz")
	c.Assert(form.Bar.Data.(string), DeepEquals, "quux")
	c.Assert(form.Baz.Data.(string), DeepEquals, "thirty-seven")
}

//Calling IsValid on a TestForm with valid data should return true
func (s *FormSuite) TestCallingIsvalidOnTestformValidDataShouldReturnTrue(c *C) {
	form := &WidgetForm{}
	Init(form)
	Populate(form, map[string]interface{}{
		"Foo":  "baz",
		"Bar":  "quux",
		"Baz":  "37",
		"Quux": "52",
	})
	c.Assert(IsValid(form), Equals, true)
}

//Calling IsValid on a TestForm with invalid data should return false
func (s *FormSuite) TestCallingIsvalidOnTestformInvalidDataShouldReturnFalse(c *C) {
	form := &WidgetForm{}
	Init(form)
	Populate(form, map[string]interface{}{
		"Foo":  "baz",
		"Bar":  "quux",
		"Baz":  "thirty-seven",
		"Quux": "52",
	})
	c.Assert(IsValid(form), Equals, false)
}

//Calling IsValid on a form with a custom Clean method should run the method on the field
func (s *FormSuite) TestCallingIsvalidOnFormCustomCleanMethodShouldRunMethodOnField(c *C) {
	form := &WidgetForm{}
	Init(form)
	Populate(form, map[string]interface{}{
		"Foo":  "baz",
		"Bar":  "quux",
		"Baz":  "37",
		"Quux": "28",
	})
	c.Assert(IsValid(form), Equals, false)
}

//Calling Copy on a valid WidgetForm and a Widget should copy the values to the widget
func (s *FormSuite) TestCallingCopyOnValidWidgetformAndWidgetShouldCopyValuesToWidget(c *C) {
	form := &WidgetForm{}
	Init(form)
	widget := &Widget{}
	Populate(form, map[string]interface{}{
		"Foo":  "baz",
		"Bar":  "quux",
		"Baz":  "37",
		"Quux": 52,
	})
	Copy(form, widget)
	c.Assert(widget.Foo, Equals, "baz")
	c.Assert(widget.Bar, Equals, "quux")
	c.Assert(widget.Baz, Equals, 37)
	c.Assert(widget.Quux, Equals, 52)
}

//Calling Copy on a valid WidgetForm and a struct missing a field from the form should return an error
func (s *FormSuite) TestCallingCopyOnValidWidgetformAndStructMissingFieldFromFormShouldReturnError(c *C){
	form := &WidgetForm{}
	Init(form)
	widget := &BrokenWidget{}
	Populate(form, map[string]interface{}{
		"Foo":  "baz",
		"Bar":  "quux",
		"Baz":  37,
		"Quux": "52",
	})
	err := Copy(form, widget)
	c.Assert(err, ErrorMatches, `No "Quux" field found on forms.BrokenWidget struct`)
}


//Calling Copy on an invalid WidgetForm and a Widget should return an error
func (s *FormSuite) TestCallingCopyOnInvalidWidgetformAndWidgetShouldCopyValuesToWidget(c *C) {
	form := &WidgetForm{}
	Init(form)
	widget := &Widget{}
	Populate(form, map[string]interface{}{
		"Foo":  "baz",
		"Bar":  "quux",
		"Baz":  "thirty-seven",
		"Quux": "52",
	})
	err := Copy(form, widget)
	c.Assert(err, ErrorMatches, "Cannot copy from invalid form")
}
