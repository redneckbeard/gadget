package gadget

import (
	"github.com/redneckbeard/gadget/routing"
	"net/http"
)

func Go(port string) {
	http.HandleFunc("/", routing.Handler())
	http.ListenAndServe(":"+port, nil)
}
