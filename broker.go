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

type contentType struct {
	mimes []string
	app   *App
}

func (ct *contentType) Via(broker Broker) {
	for _, mime := range ct.mimes {
		ct.app.Brokers[mime] = broker
	}
}

func (a *App) Accept(mimes ...string) *contentType {
	if a.Brokers == nil {
		a.Brokers = make(map[string]Broker)
	}
	return &contentType{mimes, a}
}

func (a *App) Process(r *Request, status int, body interface{}, mimetype string, data *RouteData) (int, string, string, bool) {
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
