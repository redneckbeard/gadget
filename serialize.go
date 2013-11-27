package gadget

import (
	"encoding/json"
	"encoding/xml"
	"github.com/redneckbeard/gadget/env"
)

// JsonBroker attempts to transform an interface{} value into a JSON string.
func JsonBroker(r *Request, status int, body interface{}, data *RouteData) (int, string) {
	var (
		serialized []byte
		err        error
	)
	if env.Debug {
		serialized, err = json.MarshalIndent(body, "", "  ")
	} else {
		serialized, err = json.Marshal(body)
	}
	if err != nil {
		return 500, ""
	}
	return status, string(serialized)
}

// XmlBroker attempts to transform an interface{} value into a serialized XML string.
func XmlBroker(r *Request, status int, body interface{}, data *RouteData) (int, string) {
	var (
		serialized []byte
		err        error
	)
	if env.Debug {
		serialized, err = xml.MarshalIndent(body, "", "  ")
	} else {
		serialized, err = xml.Marshal(body)
	}
	if err != nil {
		return 500, ""
	}
	return status, xml.Header + string(serialized)
}
