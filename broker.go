package gadget

import (
	"fmt"
	"net/http"
	"strings"
)

// Broker functions transform an interface{} value into a string for the body of a response
type Broker func(*Request, int, interface{}, *RouteData) (int, string)

type RouteData struct {
	ControllerName, Action, Verb string
}

var brokers = make(map[string]Broker)

type contentType struct {
	mimes []string
}

func (ct *contentType) Via(broker Broker) {
	for _, mime := range ct.mimes {
		brokers[mime] = broker
	}
}

func Accept(mimes ...string) *contentType {
	return &contentType{mimes}
}

func clearBrokers() {
	for k, _ := range brokers {
		delete(brokers, k)
	}
}

func Process(r *Request, status int, body interface{}, mimetype string, data *RouteData) (int, string, string, bool) {
	var (
		broker  Broker
		matched string
	)
	for _, mt := range strings.Split(mimetype, ",") {
		_, ok := brokers[mt]
		if ok {
			broker = brokers[mt]
			matched = mt
			break
		}
	}
	if broker == nil {
		bodyContent, ok := body.(string)
		if !ok {
			bodyContent = fmt.Sprint(body)
		}
		sniffed := http.DetectContentType([]byte(bodyContent))
		return status, bodyContent, sniffed, false
	}
	status, processed := broker(r, status, body, data)
	return status, processed, matched, true
}
