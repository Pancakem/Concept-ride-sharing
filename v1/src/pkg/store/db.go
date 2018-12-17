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

var db *sql.DB

func init() {

	ma := getConfig()
	conString := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		ma["user"].(string), ma["password"].(string), ma["dbname"].(string), ma["host"].(string),
		string(ma["port"].(int)), ma["sslmode"].(string))

	d, err := sql.Open("postgres", conString)
	if err != nil {
		log.Fatal("Couldn't establish database connection:", err)
	}
	db = d
}

func getConfig() map[string]interface{} {
	f, err := os.Open("config.yaml")
	if err != nil {
		log.Fatal("Failed to open configuration file:", err)
	}
	data, err := ioutil.ReadAll(f)
	defer f.Close()
	ma := make(map[string]interface{})
	err = yaml.Unmarshal(data, ma)
	if err != nil {
		log.Println("Couldn't unmarshal yaml data :", err)
	}
	return ma
}
