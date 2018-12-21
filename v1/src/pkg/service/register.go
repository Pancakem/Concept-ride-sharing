package service

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"

	"github.com/pancakem/rides/v1/src/pkg/common"
	"github.com/pancakem/rides/v1/src/pkg/model"
)

func RegisterUser(requestUser *model.Rider) (string, int) {
	password, _ := bcrypt.GenerateFromPassword([]byte(requestUser.Password), 5)
	requestUser.Password = string(password)
	requestUser.ID, _ = common.NewID()
	if !requestUser.Exist() {
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
	return "User already exists", 409
}

func SendMail(name, email string) {
	data := struct {
		Name string
		URL  string
	}{
		name,
		MakeUrl(email),
	}
	fmt.Println(data.URL)
	r := NewRequest([]string{email}, "Confirm Email Address", "")
	err := r.ParseTemplate("mail.html", data)
	if err != nil {
		ok, _ := r.SendMail()
		if !ok {
			r.SendMail()
		}
	}
}

func RegisterDriver(requestUser *model.Driver) (string, int) {
	fmt.Println("registration driver")
	password, _ := bcrypt.GenerateFromPassword([]byte(requestUser.Password), 5)
	requestUser.Password = string(password)
	requestUser.ID, _ = common.NewID()
	if !requestUser.Exist() {
		err := requestUser.Create()

		if err != nil {
			return err.Error(), 500
		}
		if requestUser.Email != "" {
			SendMail(requestUser.FullName, requestUser.Email)
		}

		return requestUser.ID, 201

	}
	return "User already exists", 409

}
