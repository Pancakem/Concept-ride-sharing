package auth

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pancakem/rides/v1/src/pkg/common"
	"github.com/pancakem/rides/v1/src/pkg/setting"
	"github.com/pancakem/rides/v1/src/pkg/store"
	"golang.org/x/crypto/bcrypt"
)

// JWTAuthBackend is the container for the token gen keys
type JWTAuthBackend struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

func tokenDuration() int {
	return 72
}

func expireOffset() int {
	return 3600
}

var authBackendInstance *JWTAuthBackend

// InitJWTAuthBackend initializes an auth object with a private and public key
func InitJWTAuthBackend() *JWTAuthBackend {
	if authBackendInstance == nil {
		authBackendInstance = &JWTAuthBackend{
			PrivateKey: getPrivateKey(),
			PublicKey:  getPublicKey(),
		}
	}
	return authBackendInstance
}

// GenerateToken uses the uuid of a user to give them an auth token
func (b *JWTAuthBackend) GenerateToken(userUUID string) (string, error) {
	token := jwt.New(jwt.SigningMethodRS512)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(setting.Get().JWTExpirationDelta)).Unix()
	claims["iat"] = time.Now().Unix()
	claims["sub"] = userUUID
	tokenString, err := token.SignedString(b.PrivateKey)
	return tokenString, err
}

// Authenticate validates a user
func (b *JWTAuthBackend) Authenticate(user *store.LoginForm) (bool, string, string) {
	pass := user.Password
	dbInt, err := user.Get()
	if err != nil {
		common.Log.Println(err)
		return false, "", "false"
	}

	dbRider, ok := dbInt.(store.Rider)
	if !ok {
		dbDriver, ok := dbInt.(store.Driver)
		if ok {
			if err = bcrypt.CompareHashAndPassword([]byte(dbDriver.Password), []byte(pass)); err != nil {
				return false, "", ""
			}
			return true, dbDriver.FullName, dbDriver.ID
		}
	}

	if err = bcrypt.CompareHashAndPassword([]byte(dbRider.Password), []byte(pass)); err != nil {
		return false, "", ""
	}
	return true, dbRider.FullName, dbRider.ID
}

func (b *JWTAuthBackend) getTokenRemainingValidity(timestamp interface{}) int {
	if validity, ok := timestamp.(float64); ok {
		tm := time.Unix(int64(validity), 0)
		remainder := tm.Sub(time.Now())
		if remainder > 0 {
			return int(remainder.Seconds() + float64(expireOffset()))
		}
	}
	return expireOffset()
}

// Logout expires an overdue token by adding it to the blacklist
func (b *JWTAuthBackend) Logout(tokenString string, token *jwt.Token) error {
	bl := store.BlackListed{Token: tokenString}
	return bl.Create()
}

// IsInBlacklist checks if the token has expired or was banned
func (b *JWTAuthBackend) IsInBlacklist(token string) bool {
	ok := store.GetToken(token)
	if !ok {
		return false
	}
	return true
}

func getPrivateKey() *rsa.PrivateKey {
	privateKeyFile, err := os.Open(setting.Get().PrivateKeyPath)
	if err != nil {
		common.Log.Println(err)
		return nil
	}
	permfileinfo, _ := privateKeyFile.Stat()
	var size = permfileinfo.Size()
	permbytes := make([]byte, size)

	buffer := bufio.NewReader(privateKeyFile)
	_, err = buffer.Read(permbytes)

	data, _ := pem.Decode([]byte(permbytes))

	privateKeyFile.Close()
	privateKeyImported, err := x509.ParsePKCS1PrivateKey(data.Bytes)
	if err != nil {
		common.Log.Println(err)
		return nil
	}
	return privateKeyImported
}

func getPublicKey() *rsa.PublicKey {
	publicKeyFile, err := os.Open(setting.Get().PublicKeyPath)
	if err != nil {
		common.Log.Println(err)
		return nil
	}
	permfileinfo, _ := publicKeyFile.Stat()
	var size = permfileinfo.Size()
	permbytes := make([]byte, size)

	buffer := bufio.NewReader(publicKeyFile)
	_, err = buffer.Read(permbytes)

	data, _ := pem.Decode([]byte(permbytes))

	publicKeyFile.Close()
	publicKeyImported, err := x509.ParsePKIXPublicKey(data.Bytes)
	if err != nil {
		common.Log.Println(err)
		return nil
	}

	rsaPub, _ := publicKeyImported.(*rsa.PublicKey)

	return rsaPub
}