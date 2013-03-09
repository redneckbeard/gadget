package gadget

import (
	"github.com/redneckbeard/gadget/env"
	"net/http"
)

func Go(port string) {
	err := env.Configure()
	if err != nil {
		panic(err)
	}
	env.ServeStatic()
	http.HandleFunc("/", Handler())
	http.ListenAndServe(":"+port, nil)
}
