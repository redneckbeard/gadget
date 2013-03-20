package forms

import (
	"errors"
	"fmt"
	"reflect"
)

// the Form interface has one method only required, SetOptions.
// This is where we configure how stuff is actually rendered in the form.
type Form interface {
	SetOptions()
	FieldMap() map[string]FormField
	Ready() bool
	ErrorMap() map[string]error
}

// The Default form can be embedded to provide, uh, default, behavior.
type DefaultForm struct {
	Fields map[string]FormField
	Errors map[string]error
}

func NewDefaultForm() *DefaultForm {
	return &DefaultForm{
		Fields: map[string]FormField{},
		Errors: map[string]error{},
	}
}

func (f *DefaultForm) SetOptions()                    {}
func (f *DefaultForm) Ready() bool                    { return f != nil }
func (f *DefaultForm) FieldMap() map[string]FormField { return f.Fields }
func (f *DefaultForm) ErrorMap() map[string]error   { return f.Errors }


func Init(form Form) {
	t := reflect.TypeOf(form).Elem()
	v := reflect.ValueOf(form).Elem()

	defaultFormField := v.FieldByName("DefaultForm")

	// Initialize a DefaultForm struct with a non-nil Fields map
	defaultForm := NewDefaultForm()
	defaultFormField.Set(reflect.ValueOf(defaultForm))

	fieldMap := defaultForm.Fields
	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)
		fv := v.Field(i)
		if _, ok := fv.Interface().(FormField); ok {
			// Initialize a BaseField
			base := reflect.ValueOf(NewBaseField())
			// Initialize a struct of the correct FormField type
			newField := reflect.New(fv.Type().Elem())
			newField.Elem().FieldByName("BaseField").Set(base)

			fv.Set(newField)
			fieldMap[ft.Name] = fv.Interface().(FormField)
		}
	}
}

// The Populate function copies data from a request payload into a form struct
func Populate(form Form, payload map[string]interface{}) error {
	if !form.Ready() {
		return errors.New("Cannot populate uninitialized form")
	}
	for name, field := range form.FieldMap() {
		data := payload[name]
		field.Set(data)
	}
	return nil
}

// Run a validation on the form with IsValid. With no configuration, this simply requires all values to exist and to be coercible from the value in their payload to their type.
// Optionally, Django-style `Clean<Fieldname>` methods can be specified that set the value from the payload in the form's Field object.
func IsValid(f Form) bool {
	if !f.Ready() {
		return false
	}
	t := reflect.TypeOf(f)
	for name, field := range f.FieldMap() {
		field.Clean()
		if customClean, ok := t.MethodByName("Clean" + name); ok {
			args := []reflect.Value{reflect.ValueOf(f), reflect.ValueOf(field)}
			customClean.Func.Call(args)
		}
		if field.Error() != nil {
			f.ErrorMap()[name] = field.Error()
		}
	}
	if len(f.ErrorMap()) == 0 {
		return true
	}
	return false
}

// Copy final data from form into struct
// If there's a mismatch between the exported fields of the form and the exported fields of the model struct, forms can provide a `Copy(*model)` method -- otherwise a 1:1 field match is assumed
func Copy(f Form, v interface{}) error {
	if !IsValid(f) {
		return errors.New("Cannot copy from invalid form")
	}
	structValue := reflect.ValueOf(v).Elem()
	for name, field := range f.FieldMap() {
		target := structValue.FieldByName(name)		
		if !target.IsValid() {
			return errors.New(fmt.Sprintf(`No "%s" field found on %v struct`, name, reflect.TypeOf(v).Elem()))
		}
		if !target.CanSet () {
			return errors.New(`Cannot set value on "%s" field`)
		}
		field.Copy(target)
	}
	return nil
}
