package service

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"hash"
	"strings"
	"time"

	auth "github.com/pancakem/rides/v1/src/pkg/authentication"
	"github.com/pancakem/rides/v1/src/pkg/store"
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

func decodeToken(token string) []string {
	plainText, err := rsa.DecryptOAEP(
		has,
		rand.Reader,
		keys.PrivateKey,
		[]byte(token),
		[]byte(""),
	)
	if err != nil {
		return nil
	}
	li := separate(string(plainText))
	return li
}

// MakeURL builds a URL for the email confirmation
func MakeURL(email string) string {
	ur := "/api/v1/user/confirm-email/"
	ur += generateToken(email)
	return ur
}

func separate(token string) []string {
	return strings.SplitAfter(token, " ")
}

// ValidateToken determines if a token is valid or not
func ValidateToken(token string) bool {
	li := decodeToken(token)
	if len(li) < 1 {
		return false
	}
	user := new(store.Rider)
	driv := new(store.Driver)

	user.Email = string(li[0])
	if !store.Exist(user) {
		driv.Email = string(li[0])
		if !store.Exist(driv) {
			return false
		}
		timestamp := li[1]
		if expiredTime(timestamp) {
			return false
		}
	}
	timestamp := li[1]
	if expiredTime(timestamp) {
		return false
	}
	return true
}

func expiredTime(timestamp string) bool {
	layout := "2006-01-02T15:04:05.000Z"
	t, _ := time.Parse(layout, timestamp)
	elapsed := time.Since(t)
	if elapsed > tokenLifetime {
		return false
	}
	return true
}
