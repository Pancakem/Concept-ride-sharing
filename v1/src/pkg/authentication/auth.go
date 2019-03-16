package auth

import (
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

// TokenAuthentication is the struct that packs and unpacks the tken string to json
type TokenAuthentication struct {
	Token string `json:"token" form:"token"`
}

// RequireTokenAuth ensures protected urls are accessed using an uathenticated user
func RequireTokenAuth(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	authBackend := InitJWTAuthBackend()

	token, err := jwt.Parse(req.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return authBackend.PublicKey, nil

	})

	if err == nil && token.Valid && !authBackend.IsInBlacklist(req.Header.Get("Authorization")) {
		next(res, req)
	} else {
		res.WriteHeader(http.StatusUnauthorized)
	}
}
