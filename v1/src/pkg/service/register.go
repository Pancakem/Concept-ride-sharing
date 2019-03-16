package service

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"

	"github.com/pancakem/rides/v1/src/pkg/common"
	"github.com/pancakem/rides/v1/src/pkg/store"
)

// RegisterUser creates a new user entry
func RegisterUser(user store.DefaultService) (string, int) {
	requestDriver, ok := user.(*store.Driver)
	if !ok {
		requestUser, ok := user.(*store.Rider)
		if !ok {
			return "", 0
		}
		password, _ := bcrypt.GenerateFromPassword([]byte(requestUser.Password), 5)
		requestUser.Password = string(password)
		requestUser.ID, _ = common.NewID()
		if !store.Exist(requestUser) {
			err := requestUser.Create()

			if err != nil {
				log.Println(err)
				return err.Error(), 500
			}
			if requestUser.Email != "" {
				SendMail(requestUser.FullName, requestUser.Email)
			}

			return requestUser.ID, 201
		}
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(requestDriver.Password), 5)
	requestDriver.Password = string(password)
	requestDriver.ID, _ = common.NewID()
	if !store.Exist(requestDriver) {
		err := requestDriver.Create()

		if err != nil {
			log.Println(err)
			return err.Error(), 500
		}
		if requestDriver.Email != "" {
			SendMail(requestDriver.FullName, requestDriver.Email)
		}

		return requestDriver.ID, 201
	}
	return "User already exists", 409
}

// SendMail sends an email to a given user address
func SendMail(name, email string) {
	data := struct {
		Name string
		URL  string
	}{
		name,
		MakeURL(email),
	}
	fmt.Println(data.URL)
	r := NewRequest([]string{email}, "Confirm Email Address", "")
	err := r.ParseTemplate("mail.html", data)
	if err != nil {
		err := r.SendMail()
		if err == nil {
			r.SendMail()
		}
	}
}
