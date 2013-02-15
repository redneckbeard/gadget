package controller

import "github.com/redneckbeard/gadget/requests"

// A filter takes a requests.Request and returns either nil or normal return values
type Filter func(*requests.Request) (int, interface{})
