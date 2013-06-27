package forms

import (
	"errors"
	"fmt"
	"reflect"
)

// Form is the interface that wraps the validation functionality provided by this
// package. To implement the Form interface, struct types must embed a pointer to
// DefaultForm, since the interface includes unexported fields that DefaultForm
// provides. DefaultForm implements all the methods necessary to satisfy the interface.
type Form interface {
	HasErrors() bool
	errorMap() map[string]error
	fieldMap() map[string]FormField
	ready() bool
}

// DefaultForm provides a complete implementation of the Form interface for you to
// embed in your form struct.
type DefaultForm struct {
	Fields map[string]FormField
	Errors map[string]error
}

// HasErrors returns true if the Error method of any of a Form's FormField
// members returns non-nil.
func (f *DefaultForm) HasErrors() bool { return len(f.Errors) > 0 }

func (f *DefaultForm) errorMap() map[string]error     { return f.Errors }
func (f *DefaultForm) fieldMap() map[string]FormField { return f.Fields }
func (f *DefaultForm) ready() bool                    { return f != nil }

func newDefaultForm() *DefaultForm {
	return &DefaultForm{
		Fields: make(map[string]FormField),
		Errors: make(map[string]error),
	}
}

// Init correctly initializes all FormFields and the embedded DefaultForm on a
// struct. Its first argument a Form, and its second an optional func that will
// receive the form as its only argument after initialization is complete. This
// callback enables you to do additional setup on FormFields included in your
// Form.
func Init(form Form, initializer func(Form)) {
	t := reflect.TypeOf(form).Elem()
	v := reflect.ValueOf(form).Elem()

	defaultFormField := v.FieldByName("DefaultForm")

	// Initialize a DefaultForm struct with a non-nil Fields map
	defaultForm := newDefaultForm()
	defaultFormField.Set(reflect.ValueOf(defaultForm))

	fieldMap := defaultForm.Fields
	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)
		fv := v.Field(i)
		if _, ok := fv.Interface().(FormField); ok {
			// Initialize a BaseField
			base := reflect.ValueOf(newBaseField())
			// Initialize a struct of the correct FormField type
			newField := reflect.New(fv.Type().Elem())
			newField.Elem().FieldByName("BaseField").Set(base)

			fv.Set(newField)
			formField := fv.Interface().(FormField)
			formField.DefaultMessages()
			fieldMap[ft.Name] = formField
		}
	}
	if initializer != nil {
		initializer(form)
	}
}

// Populate copies data from a map[string]interface{} value into the Data members
// of the corresponding FormFields of the Form. The keys of the map must be a
// case-sensitive exact match of exported names of FormFields on the struct.
// Populate will return an error if Init has not yet been called on the Form.
// 
// Note that the type of the second argument is not map[string][]string. This is
// meant to provide ease of use with gadget.Request.Params over http.Request.Form.
func Populate(form Form, payload map[string]interface{}) error {
	if !form.ready() {
		return errors.New("Cannot populate uninitialized form")
	}
	for name, field := range form.fieldMap() {
		data := payload[name]
		field.Set(data)
	}
	return nil
}

// IsValid calls the Clean method on all of a Form's FormFields. As a result, if
// IsValid returns true, the Value fields of all required FormFields will be
// non-nil. If false, error messages may be retrieved by the Error method of an
// individual field or by retrieving the error for a field by name from the Form's
// Errors map.
// 
// If validation is needed beyond what a FormField type provides by default, Forms
// may define special methods to perform more robust checks on the value of the
// Data field. These methods are named "Clean<Fieldname>" and receive the field
// being validated as their argument. An example is provided below.
func IsValid(f Form) bool {
	if !f.ready() {
		return false
	}
	t := reflect.TypeOf(f)
	for name, field := range f.fieldMap() {
		if field.isNil() {
			if field.isRequired() {
				field.SetError("This field is required")
			}
		} else {
			field.Clean()
			if customClean, ok := t.MethodByName("Clean" + name); ok {
				args := []reflect.Value{reflect.ValueOf(f), reflect.ValueOf(field)}
				customClean.Func.Call(args)
			}
		}
		if field.Error() != nil {
			f.errorMap()[name] = field.Error()
		}
	}
	if len(f.errorMap()) == 0 {
		return true
	}
	return false
}

// Copy transfers data from a validated Form to another struct for further
// processing (presumably persistance). The target struct must have a field of
// the appropriate type and the same name for every FormField in the Form.
//
// 
// Copy will return an error if the Form has not passed validation, if it cannot
// find a corresponding field on the target struct for a FormField, or if it
// cannot set the FormField's value on the target field.
func Copy(f Form, target interface{}) error {
	if !IsValid(f) {
		return errors.New("Cannot copy from invalid form")
	}
	structValue := reflect.ValueOf(target).Elem()
	for name, field := range f.fieldMap() {
		if field.canCopy() {
			targetField := structValue.FieldByName(name)
			if !targetField.IsValid() {
				return errors.New(fmt.Sprintf(`No "%s" field found on %v struct`, name, reflect.TypeOf(target).Elem()))
			}
			if !targetField.CanSet() {
				return errors.New(`Cannot set value on "%s" field`)
			}
			field.Copy(targetField)
		}
	}
	return nil
}
