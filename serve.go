package gadget

import (
	"github.com/redneckbeard/gadget/env"
	"github.com/redneckbeard/gadget/routing"
	"net/http"
)

func Go(port string) {
	env.Configure()
	env.ServeStatic()
	http.HandleFunc("/", routing.Handler())
	http.ListenAndServe(":"+port, nil)
}
