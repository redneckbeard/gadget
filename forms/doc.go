/*
Package forms provides tools for validating data found in
gadget.Request.Params. Its structure is inspired by Django's forms library and
is intended to be both declarative and highly independent of your application's
model layer. Forms are struct types that embed *DefaultForm and declare as
their fields pointers to any number of Field types. They look something like
this:

	type InvoiceForm struct {
		*DefaultForm
		Client, Notes *StringField
		Date          *TimeField
		Number        *IntField
	}

*/
package forms
