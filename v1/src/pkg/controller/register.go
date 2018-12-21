package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pancakem/rides/v1/src/pkg/authentication"

	"github.com/gorilla/mux"

	"github.com/pancakem/rides/v1/src/pkg/common"
	"github.com/pancakem/rides/v1/src/pkg/model"
	"github.com/pancakem/rides/v1/src/pkg/service"
)

// RegisterRider route
func RegisterRider(w http.ResponseWriter, r *http.Request) {
	form := model.Rider{}
	data, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	common.CheckError(err)
	err = json.Unmarshal(data, &form)
	common.CheckError(err)
	id, status := service.RegisterUser(&form)
	if status != 201 {
		w.WriteHeader(status)
		return
	}

	w.WriteHeader(status)

	authBackend := auth.InitJWTAuthBackend()
	token, err := authBackend.GenerateToken(id)
	if err != nil {
		log.Println(err)
		return
	}
	m := make(map[string]string)
	m["id"] = id
	m["name"] = form.FullName
	m["token"] = token
	data, err = json.Marshal(m)
	w.Write(data)

}

func RegisterDriver(w http.ResponseWriter, r *http.Request) {
	form := model.Driver{}
	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
	}
	id, status := service.RegisterDriver(&form)
	fmt.Println(id, status)
	if status != 201 {
		w.WriteHeader(status)
		w.Write([]byte(id))
		return
	}

	authBackend := auth.InitJWTAuthBackend()
	token, err := authBackend.GenerateToken(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(200)
		return
	}
	w.WriteHeader(status)
	form.GetByID()
	m := make(map[string]string)
	m["id"] = id
	m["name"] = form.FullName
	m["token"] = token
	m["image_url"] = form.ImageURL
	data, err := json.Marshal(m)
	w.Write(data)
}

func ConfirmEmail(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	token := params["token"]

	err := service.DecodeToken(token)

	defer func() {
		if err != nil {
			http.Error(w, "Invalid token. You have been sent a new token. Check your email", 401)
		}
	}()

}
