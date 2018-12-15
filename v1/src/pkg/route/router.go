package route

import (
	"github.com/gorilla/mux"
)

const apiVersion = "/api/v1"

// InitRouter adds all endpoints and return router
func InitRouter() *mux.Router {
	router := mux.NewRouter()
	router = GetARide(router)
	return router
}
