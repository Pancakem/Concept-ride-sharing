package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/pancakem/rides/v1/src/pkg/route"
)

var addr = flag.String("addr", "0.0.0.0:4000", "http service address")

func main() {
	flag.Parse()
	router := route.InitRouter()
	log.Println("Listening at port 4000")
	log.Fatal(http.ListenAndServe(*addr, router))
}
