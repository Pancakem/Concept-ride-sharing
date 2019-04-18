package setting

import (
	"io/ioutil"
	"os"

	"github.com/pancakem/rides/v1/src/pkg/common"
	"gopkg.in/yaml.v2"
)

var environments = map[string]string{
	"production":  "prod.yaml",
	"development": "pre.yaml",
	"tests":       "tests.yaml",
}

// Settings holds the data required to successfully start the app
type Settings struct {
	Environment        string
	PrivateKeyPath     string `yaml:"PrivateKeyPath"`
	PublicKeyPath      string `yaml:"PublicKeyPath"`
	JWTExpirationDelta int    `yaml:"JWTExpirationDelta"`
}

var settings = Settings{}

func init() {
	env := os.Getenv("GO_ENV")
	if env == "" {
		common.Log.Println("Warning: Setting development environment due to lack of GO_ENV value")
		env = "development"
	}
	LoadSettingsByEnv(env)
}

// LoadSettingsByEnv returns a pointer to a settings struct
func LoadSettingsByEnv(env string) *Settings {
	content, err := ioutil.ReadFile(environments[env])
	if err != nil {
		common.Log.Println("Error while reading config file", err)
		return nil
	}
	settings := new(Settings)
	yamlErr := yaml.Unmarshal(content, settings)
	if yamlErr != nil {
		common.Log.Println("Error while parsing config file", yamlErr)
		return nil
	}
	return nil
}

// Get returns the settings
func Get() Settings {
	return settings
}
