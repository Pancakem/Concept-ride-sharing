package main

import (
	"net/http"

	"github.com/pancakem/rides/v1/src/pkg/route"
)

func main() {
	router := route.InitRouter()
	http.ListenAndServe(":4000", router)
}
