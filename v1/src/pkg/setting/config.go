package setting

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// Config  variables
type Config struct {
	// the key for accessing google maps api
	MapKey string `yaml:"mapKey"`
	// the database host
	Host         string `yaml:"host"`
	Port         int `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	DatabaseName string `yaml:"dbname"`
	SSLMode      string `yaml:"sslmode"`
	// the log url 
	LogFile string `yaml:"logurl"`
}

// Get reads the config file
func (c *Config) Get() error {
	f, err := os.Open("config.yaml")
	defer f.Close()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		log.Println("Failed to read yaml file. Exiting")
		return err
	}
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		log.Println(err)
	}
	return err
}
