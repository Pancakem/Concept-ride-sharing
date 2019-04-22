package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/pancakem/rides/v1/src/pkg/route"
)

func main() {
	var addr = flag.String("addr", "0.0.0.0:4000", "http service address")
	flag.Parse()
	router := route.InitRouter()

	log.Println("Listening at port 4000")
	log.Fatal(http.ListenAndServe(*addr, router))
}
