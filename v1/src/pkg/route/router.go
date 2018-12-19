package route

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	auth "github.com/pancakem/user-service/v1/src/pkg/authentication"
	"github.com/pancakem/user-service/v1/src/pkg/controller"
)

const apiVersion = "/api/v1"

// InitRouter adds all endpoints and return router
func InitRouter() *mux.Router {
	router := mux.NewRouter()
	router = SetAuthLayer(router)
	router = WorkUser(router)
	router = GetARide(router)
	return router
}

// WorkUser avails the
func WorkUser(router *mux.Router) *mux.Router {
	router.Handle(apiVersion+"/riders",
		negroni.New(
			negroni.HandlerFunc(auth.RequireTokenAuth),
			negroni.HandlerFunc(controller.GetRiders),
		)).Methods("GET")
	return router
}

// SetAuthLayer returns authorization and registration router
func SetAuthLayer(router *mux.Router) *mux.Router {
	router.HandleFunc(apiVersion+"/user/confirm-email/{token}", controller.ConfirmEmail).Methods("GET")
	router.HandleFunc(apiVersion+"/rider/register", controller.RegisterRider).Methods("POST")
	router.HandleFunc(apiVersion+"/driver/register", controller.RegisterDriver).Methods("POST")
	router.HandleFunc(apiVersion+"/login", controller.Login).Methods("POST")
	router.Handle(apiVersion+"/refresh-token-auth", negroni.New(
		negroni.HandlerFunc(auth.RequireTokenAuth),
		negroni.HandlerFunc(controller.RefreshToken),
	)).Methods("GET")
	router.Handle(apiVersion+"/logout",
		negroni.New(
			negroni.HandlerFunc(auth.RequireTokenAuth),
			negroni.HandlerFunc(controller.Logout),
		)).Methods("GET")
	return router
}
