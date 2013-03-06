package processor

import (
	"fmt"
	"net/http"
	"strings"
)

// Processor functions transform an interface{} value into a string for the body of a response
type Processor func(int, interface{}, *RouteData) (int, string)

type RouteData struct {
	ControllerName, Action, Verb string
}

var processors = make(map[string]Processor)

func init() {
}

func Define(mimetype string, processor Processor) {
	processors[mimetype] = processor
}

func clear() {
	for k, _ := range processors {
		delete(processors, k)
	}
}

func Process(status int, body interface{}, mimetype string, data *RouteData) (int, string, string, bool) {
	var (
		processor Processor
		matched string
	)
	for _, mt := range strings.Split(mimetype, ",") {
		_, ok := processors[mt]
		if ok {
			processor = processors[mt]
			matched = mt
			break
		}
	}
	if processor == nil {
		bodyContent, ok := body.(string)
		if !ok {
			bodyContent = fmt.Sprint(body)
		}
		sniffed := http.DetectContentType([]byte(bodyContent))
		return status, bodyContent, sniffed, false
	}
	status, processed := processor(status, body, data)
	return status, processed, matched, true
}
