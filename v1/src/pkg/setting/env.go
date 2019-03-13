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

// Settings holds the data required to successfully start the app
type Settings struct {
	Environment        string
	PrivateKeyPath     string `josn:"PrivateKeyPath"`
	PublicKeyPath      string `json:"PublicKeyPath"`
	JWTExpirationDelta int    `json:"JWTExpirationDelta"`
}

var settings = Settings{}

func init() {
	env := os.Getenv("GO_ENV")
	if env == "" {
		log.Println("Warning: Setting development environment due to lack of GO_ENV value")
		env = "development"
	}
	LoadSettingsByEnv(env)
}

// LoadSettingsByEnv returns a pointer to a settings struct
func LoadSettingsByEnv(env string) *Settings {
	content, err := ioutil.ReadFile(environments[env])
	if err != nil {
		log.Println("Error while reading config file", err)
		return nil
	}
	settings := new(Settings)
	jsonErr := json.Unmarshal(content, settings)
	if jsonErr != nil {
		log.Println("Error while parsing config file", jsonErr)
		return nil
	}
	return nil
}

// Get returns the settings
func Get() Settings {
	return settings
}
