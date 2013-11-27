package gadget

import (
	"fmt"
	"net/http"
	"strings"
)

// Broker functions transform an interface{} value into a string for the body of a response
type Broker func(*Request, int, interface{}, *RouteData) (int, string)

// RouteData provides context about the where the current request is being executed.
type RouteData struct {
	ControllerName, Action, Verb string
}

type contentType struct {
	mimes []string
	app   *App
}

// Via associates a Broker with a MIME type registered with Accept.
func (ct *contentType) Via(broker Broker) {
	for _, mime := range ct.mimes {
		ct.app.Brokers[mime] = broker
	}
}

// Accept registers MIME type strings with the App. Via can be called on its return value
// to associated those MIME types with a Broker.
func (a *App) Accept(mimes ...string) *contentType {
	if a.Brokers == nil {
		a.Brokers = make(map[string]Broker)
	}
	return &contentType{mimes, a}
}

func (a *App) process(r *Request, status int, body interface{}, mimetype string, data *RouteData) (int, string, string, bool) {
	var (
		broker  Broker
		matched string
	)
	for _, mt := range strings.Split(mimetype, ",") {
		_, ok := a.Brokers[mt]
		if ok {
			broker = a.Brokers[mt]
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
