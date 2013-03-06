package processor

import (
	"encoding/json"
	"encoding/xml"
)

func JsonProcessor(status int, body interface{}, action string) (int, string) {
	serialized, err := json.Marshal(body)
	if err != nil {
		return 500, ""
	}
	return status, string(serialized)
}

func XmlProcessor(status int, body interface{}, action string) (int, string) {
	serialized, err := xml.MarshalIndent(body, "", " ")
	if err != nil {
		return 500, ""
	}
	return status, string(serialized)
}
