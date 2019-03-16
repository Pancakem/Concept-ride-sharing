package controller

import (
	"encoding/json"
	"net/http"

	"github.com/pancakem/rides/v1/src/pkg/service"
	"github.com/pancakem/rides/v1/src/pkg/store"
)

// Login handles a auth request
func Login(w http.ResponseWriter, r *http.Request) {
	requestUser := new(store.Rider)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestUser)
	lf := &store.LoginForm{Email: requestUser.Email, Phonenumber: requestUser.Phonenumber, Password: requestUser.Password}

	responseStatus, token := service.Login(lf)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseStatus)
	w.Write(token)
}

// RefreshToken handles a token refresh request
func RefreshToken(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	requestUser := new(store.Rider)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestUser)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", string(service.RefreshToken(requestUser)))
}

// Logout handles a logout request
func Logout(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	err := service.Logout(r)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
