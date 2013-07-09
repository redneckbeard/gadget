package forms

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// FormField is the interface that encapsulates field-level form operations.
// All types implementing FormField should be struct types that embed
// BaseField, which provides all methods except Clean and Copy. Individual
// field types must implement those methods appropriatly for the types of data
// they are intended to validate.
//
// By virtue of their embedded BaseField, all FormFields have a Data field of
// type interface{}. The FormFields themselves conventionally have a Value
// field of the type that they purport to validate. All FormFields are required
// unless their Required field is set to false in the form's SetOptions method.
//
// The forms package ships with FormField types for string, int, bool, float64,
// and time.Time values. Users can easily define additional fields. Refer to
// the source for the builtin FormField types for details of how to define new
// ones.
type FormField interface {
	Clean()
	Copy(reflect.Value)
	canCopy() bool
	DefaultMessages()
	Error() error
	GetMessage(string) string
	Set(interface{})
	isNil() bool
	SetError(string)
	SetMessage(string, string)
	isRequired() bool
}

type BaseField struct {
	err         error
	Data        interface{}
	Messages    map[string]string
	Required    bool
	Placeholder bool
}

func newBaseField() *BaseField {
	messages := make(map[string]string)
	return &BaseField{Messages: messages, Required: true}
}

func (field *BaseField) isNil() bool {
	return field.Data == nil
}

func (field *BaseField) isRequired() bool {
	return field.Required
}

func (field *BaseField) canCopy() bool {
	return !field.Placeholder
}

func (field *BaseField) Set(data interface{}) {
	field.Data = data
}

func (field *BaseField) SetError(msg string) {
	field.err = errors.New(msg)
}

func (field *BaseField) Error() error {
	return field.err
}

func (field *BaseField) GetMessage(msgName string) string {
	return field.Messages[msgName]
}

func (field *BaseField) SetMessage(msgName string, msgValue string) {
	field.Messages[msgName] = msgValue
}

func (field *BaseField) DefaultMessages() {}

// StringField validates the presences of a string value.
type StringField struct {
	*BaseField
	Value string
}

func (field *StringField) Clean() {
	switch field.Data.(type) {
	case string:
		field.Value = field.Data.(string)
	default:
		field.SetError(field.GetMessage("type"))
	}
}

func (field *StringField) Copy(v reflect.Value) {
	v.SetString(field.Value)
}

func (field *StringField) DefaultMessages() {
	field.SetMessage("type", "A string value is required")
}

// IntField validates the presence of an int value. If the Data on the
// FormField is string, its Clean method will attempt to parse it.
type IntField struct {
	*BaseField
	Value int64
}

func (field *IntField) Clean() {
	switch field.Data.(type) {
	case int, int64:
		field.Value = int64(field.Data.(int))
	case float64:
		n := field.Data.(float64)
		if n - float64(int(n)) == 0.0 {
			field.Value = int64(n)
		} else {
			field.SetError(field.GetMessage("type"))
		}
	case string:
		intString := field.Data.(string)
		i, err := strconv.ParseInt(intString, 10, 0)
		if err != nil {
			field.SetError(field.GetMessage("type"))
		} else {
			field.Value = i
		}
	default:
		field.SetError(field.GetMessage("type"))
	}
}

func (field *IntField) Copy(v reflect.Value) {
	v.SetInt(int64(field.Value))
}

func (field *IntField) DefaultMessages() {
	field.SetMessage("type", "An integer value is required")
}

type IntSliceField struct {
	*BaseField
	Value []int64
}

func (field *IntSliceField) Clean() {
	switch field.Data.(type) {
	case int64:
		field.Value = []int64{field.Data.(int64)}
	case string:
		s := field.Data.(string)
		i, err := strconv.ParseInt(s, 10, 0)
		if err != nil {
			field.SetError(field.GetMessage("type"))
			return
		} else {
			field.Value = []int64{i}
		}
	case []int64:
		field.Value = field.Data.([]int64)
	case []string:
		intStrings := field.Data.([]string)
		ints := []int64{}
		for _, s := range intStrings {
			i, err := strconv.ParseInt(s, 10, 0)
			if err != nil {
				field.SetError(field.GetMessage("type"))
				return
			} else {
				ints = append(ints, i)
			}
		}
		field.Value = ints
	default:
		field.SetError(field.GetMessage("type"))
	}
}

func (field *IntSliceField) Copy(v reflect.Value) {
	v.Set(reflect.ValueOf(field.Value))
}

func (field *IntSliceField) DefaultMessages() {
	field.SetMessage("type", "An array of integers is required")
}

// Float64Field validates the presence of a float64 value. If the Data on the
// FormField is a string, the Clean method will attempt to parse it. If it is
// an int, it will cast it to a float64.
type Float64Field struct {
	*BaseField
	Value float64
}

func (field *Float64Field) Clean() {
	switch field.Data.(type) {
	case float64:
		field.Value = field.Data.(float64)
	case int:
		field.Value = float64(field.Data.(int))
	case string:
		floatString := field.Data.(string)
		f64, err := strconv.ParseFloat(floatString, 64)
		if err != nil {
			field.SetError(field.GetMessage("type"))
		}
		field.Value = f64
	default:
		field.SetError(field.GetMessage("type"))
	}
}

func (field *Float64Field) Copy(v reflect.Value) {
	v.SetFloat(field.Value)
}

func (field *Float64Field) DefaultMessages() {
	field.SetMessage("type", "An float value is required")
}

type BoolField struct {
	*BaseField
	Value bool
}

// Boolfield validates the presence of a boolean value. If the Data on the
// FormField is an int, 0 is taken for false and 1 for true. If it is a string,
// it is parsed by strconv.ParseBool.
func (field *BoolField) Clean() {
	switch field.Data.(type) {
	case bool:
		field.Value = field.Data.(bool)
	case int:
		switch field.Data.(int) {
		case 0:
			field.Value = false
		case 1:
			field.Value = true
		default:
			field.SetError(field.GetMessage("integer_out_of_bounds"))
		}
		field.Value = field.Data.(int) != 0
	case string:
		str := field.Data.(string)
		if str == "" {
			field.Value = false
			return
		}
		value, err := strconv.ParseBool(field.Data.(string))
		if err != nil {
			field.SetError(field.GetMessage("type"))
		} else {
			field.Value = value
		}
	default:
		field.SetError(field.GetMessage("type"))
	}
}

func (field *BoolField) Copy(v reflect.Value) {
	v.SetBool(field.Value)
}

func (field *BoolField) DefaultMessages() {
	field.SetMessage("type", "An boolean value is required")
	field.SetMessage("integer_out_of_bounds", "Only 0 and 1 are valid integer values for a boolean field")
}

// TimeField validates the presence of a Time value. If the value is an int64, it
// is parsed as seconds since the epoch. If it is a string, TimeField will attempt
// to parse it with the value of its Format field as the layout. The default time
// format is RFC 33339.
type TimeField struct {
	*BaseField
	Value  time.Time
	Format string
}

func (field *TimeField) Clean() {
	if field.Format == "" {
		field.Format = time.RFC3339
	}
	typeErr := fmt.Sprintf(field.GetMessage("type"), field.Format)
	switch field.Data.(type) {
	case time.Time:
		field.Value = field.Data.(time.Time)
	case int64:
		i := field.Data.(int64)
		if i < 0 {
			field.SetError(typeErr)
		} else {
			field.Value = time.Unix(int64(i), 0).UTC()
		}
	case string:
		t, err := time.Parse(field.Format, field.Data.(string))
		if err != nil {
			field.SetError(typeErr)
		} else {
			field.Value = t
		}
	default:
		field.SetError(typeErr)
	}
}

func (field *TimeField) Copy(v reflect.Value) {
	v.Set(reflect.ValueOf(field.Value))
}

func (field *TimeField) DefaultMessages() {
	field.SetMessage("type", `An integer of seconds since the epoch or a string in the format "%s" is required`)
}
