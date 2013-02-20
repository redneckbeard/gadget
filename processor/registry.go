package processor

import "fmt"

// Processor functions transform an interface{} value into a string for the body of a response
type Processor func(int, interface{}) (int, string)

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

func Process(status int, body interface{}, mimetype string) (int, string, bool) {
	switch body.(type) {
	case string:
		return status, body.(string), false
	}
	processor, ok := processors[mimetype]
	if !ok {
		return status, fmt.Sprint(body), false
	}
	status, processed := processor(status, body)
	return status, processed, true
}
