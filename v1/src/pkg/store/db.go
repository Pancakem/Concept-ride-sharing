package store

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"

	_ "github.com/lib/pq"
)

type config struct {
	user     string `yaml:"user"`
	password string `yaml:"password"`
	dbname   string `yaml:"dbname"`
	host     string `yaml:"dbname"`
	port     int    `yaml:"port"`
	sslmode  string `yaml:"sslmode"`
}

var db *sql.DB

func init() {

	ma := getConfig()
	conString := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		ma.user, ma.password, ma.dbname, ma.host,
		string(ma.port), ma.sslmode)

	d, err := sql.Open("postgres", conString)
	if err != nil {
		log.Fatal("Couldn't establish database connection:", err)
	}
	db = d
}

func getConfig() *config {
	f, err := os.Open("config.yaml")
	if err != nil {
		log.Fatal("Failed to open configuration file:", err)
	}
	data, err := ioutil.ReadAll(f)
	defer f.Close()
	ma := &config{}
	err = yaml.Unmarshal(data, ma)
	if err != nil {
		log.Println("Couldn't unmarshal yaml data :", err)
	}
	return ma
}
