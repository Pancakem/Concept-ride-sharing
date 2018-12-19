package service

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"hash"
	"strings"
	"time"

	auth "github.com/pancakem/rides/v1/src/pkg/authentication"
	"github.com/pancakem/rides/v1/src/pkg/model"
)

const tokenLifetime = 2e15

var (
	has  hash.Hash
	keys *auth.JWTAuthBackend
)

func generateToken(user string) string {
	keys = auth.InitJWTAuthBackend()
	message := []byte(user + " " + string(time.Now().String()))
	label := []byte("")
	has = sha256.New()

	token, err := rsa.EncryptOAEP(
		has,
		rand.Reader,
		keys.PublicKey,
		message,
		label,
	)
	if err != nil {
		return ""
	}

	return string(token)
}

func DecodeToken(token string) error {
	user := &model.Rider{}
	driv := &model.Driver{}
	plainText, err := rsa.DecryptOAEP(
		has,
		rand.Reader,
		keys.PrivateKey,
		[]byte(token),
		[]byte(""),
	)
	if err != nil {
		li := separate(string(plainText))
		user.Email = string(li[0])
		if !user.Exist() {
			driv.Email = string(li[0])
			if !driv.Exist() {
				return err
			}
			timestamp := li[0]
			layout := "2006-01-02T15:04:05.000Z"
			t, _ := time.Parse(layout, timestamp)
			elapsed := time.Since(t)
			if elapsed > tokenLifetime {
				SendMail(driv.FullName, driv.Email)
				var err error
				return err
			}
			driv.IsActive = true
			return err

		}
		timestamp := li[1]
		//
		layout := "2006-01-02T15:04:05.000Z"
		t, _ := time.Parse(layout, timestamp)
		elapsed := time.Since(t)
		if elapsed > tokenLifetime {
			SendMail(user.FullName, user.Email)
			var err error
			return err
		}
		user.IsActive = true
		return err
	}
	return err
}

func MakeUrl(email string) string {
	ur := "/api/v1/user/confirm-email/"
	ur += generateToken(email)
	return ur
}

func separate(token string) []string {
	return strings.SplitAfter(token, " ")
}
