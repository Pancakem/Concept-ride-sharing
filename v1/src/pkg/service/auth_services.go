package service

import (
	"encoding/json"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"

	auth "github.com/pancakem/rides/v1/src/pkg/authentication"
	"github.com/pancakem/rides/v1/src/pkg/model"
)

// Login authenticates a user
func Login(requestUser *model.LoginForm) (int, []byte) {
	authBackend := auth.InitJWTAuthBackend()
	ok, name, id := authBackend.Authenticate(requestUser)
	if ok {
		token, err := authBackend.GenerateToken(id)
		if err != nil {
			return http.StatusInternalServerError, []byte("")
		}
		ma := make(map[string]string)
		ma["id"] = id
		ma["user"] = name
		ma["token"] = token
		ma["image_url"] = "" // add image url

		response, _ := json.Marshal(ma)
		return http.StatusOK, response
	}
	return http.StatusUnauthorized, []byte("")
}

// RefreshToken returns a token for the user given
func RefreshToken(v interface{}) []byte {
	authBackend := auth.InitJWTAuthBackend()
	requestRider, ok := v.(model.Rider)
	var token string
	var err error
	if ok {
		token, err = authBackend.GenerateToken(string(requestRider.ID))
		if err != nil {
			return []byte("")
		}
	} else {
		requestDriver := v.(model.Driver)
		token, err = authBackend.GenerateToken(string(requestDriver.ID))
		if err != nil {
			return []byte("")
		}

	}
	response, err := json.Marshal(auth.TokenAuthentication{Token: token})
	if err != nil {
		return []byte("")
	}
	return response
}

// Logout deauths a user
func Logout(req *http.Request) error {
	authBackend := auth.InitJWTAuthBackend()
	tokenRequest, err := request.ParseFromRequest(req, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
		return authBackend.PublicKey, nil
	})
	if err != nil {
		return err
	}
	tokenString := req.Header.Get("Authorization")
	return authBackend.Logout(tokenString, tokenRequest)
}
