package main

import (
	"log"
	"net/http"

	"github.com/pancakem/rides/v1/src/pkg/route"
)

func main() {
	router := route.InitRouter()
	log.Println("Listening at port 4000")
	log.Fatal(http.ListenAndServe(":4000", router))

}
