package forms_test

import (
	"fmt"
	"github.com/redneckbeard/gadget/forms"
	"regexp"
)

type UserForm struct {
	*forms.DefaultForm
	Username, Email *forms.StringField
}

func (form *UserForm) CleanEmail(f forms.FormField) {
	field := f.(*forms.StringField)
	naiveEmailPattern := regexp.MustCompile(`[\w.-]+@[\w.-]\.(com|net|org)`)
	if !naiveEmailPattern.MatchString(field.Value) {
		field.SetError("Email must be an email")
	}
}

func ExampleIsValid() {
	/*
	 *type UserForm struct {
	 *        *forms.DefaultForm
	 *        Username, Email *forms.StringField
	 *}
	 *
	 *func (form *UserForm) CleanEmail(f forms.FormField) {
	 *        field := f.(*forms.StringField)
	 *        naiveEmailPattern := regexp.MustCompile(`[\w.-]+@[\w.-]\.(com|net|org)`)
	 *        if !naiveEmailPattern.MatchString(field.Value) {
	 *                field.SetError("Email must be an email")
	 *        }
	 *}
	 */

	form := &UserForm{}
	forms.Init(form, nil)
	forms.Populate(form, map[string]interface{}{
		"Username": "penny",
		"Email":    "clearly not an email address",
	})
	fmt.Println(forms.IsValid(form))
	fmt.Println(form.Email.Error())
	// Output:
	// false
	// Email must be an email
}

func ExampleInit() {
	type ArticleForm struct {
		*forms.DefaultForm
		Title, Body *forms.StringField
		AuthorId    *forms.IntField
	}
	form := &ArticleForm{}
	forms.Init(form, func(f forms.Form) {
		form := f.(*ArticleForm)
		form.AuthorId.Required = false
	})
	forms.Populate(form, map[string]interface{}{
		"Title": "Mission: Stop Dr. Claw",
	})
	fmt.Println(forms.IsValid(form))
	fmt.Println(form.AuthorId.Error())
	fmt.Println(form.Body.Error())
	// Output:
	// false
	// <nil>
	// This field is required
}

func ExampleCopy() {
	type Invoice struct {
		Client, Notes string
		Paid          bool
		Number        int
	}

	type InvoiceForm struct {
		*forms.DefaultForm
		Client, Notes *forms.StringField
		Paid          *forms.BoolField
		Number        *forms.IntField
	}

	form := &InvoiceForm{}
	forms.Init(form, nil)
	forms.Populate(form, map[string]interface{}{
		"Client": "M.A.D.",
		"Notes":  "Some day, Gadget, some day...",
		"Paid":   "1",
		"Number": 23,
	})
	fmt.Println(forms.IsValid(form))

	invoice := &Invoice{}
	forms.Copy(form, invoice)
	fmt.Println(invoice.Client)
	fmt.Println(invoice.Notes)
	fmt.Println(invoice.Paid)
	fmt.Println(invoice.Number)
	// Output:
	// true
	// M.A.D.
	// Some day, Gadget, some day...
	// true
	// 23
}
