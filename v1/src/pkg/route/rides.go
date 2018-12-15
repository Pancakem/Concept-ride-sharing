package route

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"

	"github.com/pancakem/rides/v1/src/pkg/controller"
)

// GetARide routers
func GetARide(router *mux.Router) *mux.Router {
	router.Handle(apiVersion+"/rides/bookride",
		negroni.New(
			// negroni.HandlerFunc(auth.RequireTokenAuth),
			negroni.HandlerFunc(controller.BookRide),
		))

	router.Handle(apiVersion+"/rides/driverlocation",
		negroni.New(
			// negroni.HandlerFunc(auth.RequireTokenAuth),
			negroni.HandlerFunc(controller.GetLocation),
		))

	return router
}
