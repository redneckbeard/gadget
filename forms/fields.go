package forms

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

type FormField interface {
	Set(interface{})
	Clean()
	SetError(string)
	Error() error
	GetMessage(string) string
	SetMessage(string, string)
	DefaultMessages()
	Copy(reflect.Value)
}

type BaseField struct {
	err      error
	Data     interface{}
	Messages map[string]string
}

func NewBaseField() *BaseField {
	messages := make(map[string]string)
	return &BaseField{Messages: messages}
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

type IntField struct {
	*BaseField
	Value int
}

func (field *IntField) Clean() {
	switch field.Data.(type) {
	case int:
		field.Value = field.Data.(int)
	case string:
		intString := field.Data.(string)
		i, err := strconv.ParseInt(intString, 10, 0)
		if err != nil {
			field.SetError(field.GetMessage("type"))
		} else {
			field.Value = int(i)
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

func (field *BoolField) Clean() {
	switch field.Data.(type) {
	case bool:
		field.Value = field.Data.(bool)
	case int:
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
}

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
