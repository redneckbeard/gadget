package strutil

import (
	"testing"
)

func TestSnakify(t *testing.T) {
	cases := []struct {
		in, out string
	}{

		{in: "Product", out: "product"},
		{in: "SpecialGuest", out: "special_guest"},
		{in: "ApplicationController", out: "application_controller"},
		{in: "Area51Controller", out: "area51_controller"},
		{in: "HTMLTidy", out: "html_tidy"},
		{in: "HTMLTidyGenerator", out: "html_tidy_generator"},
		{in: "FreeBSD", out: "free_bsd"},
		{in: "HTML", out: "html"},
		{in: "ForceXMLController", out: "force_xml_controller"},
	}
	for _, c := range cases {
		if Snakify(c.in) != c.out {
			t.Fail()
			t.Log(c.in, "=>", Snakify(c.in))
		}
	}
}
