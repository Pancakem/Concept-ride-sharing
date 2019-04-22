package route

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	auth "github.com/pancakem/rides/v1/src/pkg/authentication"
	"github.com/pancakem/rides/v1/src/pkg/controller"
)

// WorkUser avails the
func WorkUser(router *mux.Router) *mux.Router {
	router.Handle(apiVersion()+"/riders",
		negroni.New(
			negroni.HandlerFunc(auth.RequireTokenAuth),
			negroni.HandlerFunc(controller.GetRiders),
		)).Methods("GET")
	return router
}
