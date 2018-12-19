package setting

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var environments = map[string]string{
	"production":  "prod.json",
	"development": "pre.json",
	"tests":       "tests.json",
}

type Settings struct {
	PrivateKeyPath     string `josn:"PrivateKeyPath"`
	PublicKeyPath      string `json:"PublicKeyPath"`
	JWTExpirationDelta int    `json:"JWTExpirationDelta"`
}

var settings = Settings{}
var env = "development"

func init() {
	env = os.Getenv("GO_ENV")
	if env == "" {
		log.Println("Warning: Setting development environment due to lack of GO_ENV value")
		env = "development"
	}
	LoadSettingsByEnv(env)
}

func LoadSettingsByEnv(env string) {
	content, err := ioutil.ReadFile(environments[env])
	if err != nil {
		log.Println("Error while reading config file", err)
	}
	settings = Settings{}
	jsonErr := json.Unmarshal(content, &settings)
	if jsonErr != nil {
		log.Println("Error while parsing config file", jsonErr)
	}
}

func GetEnvironment() string {
	return env
}

func Get() Settings {
	return settings
}

func IsTestEnvironment() bool {
	return env == "tests"
}
