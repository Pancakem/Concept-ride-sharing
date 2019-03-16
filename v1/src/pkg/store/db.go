package store

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"

	_ "github.com/lib/pq" //
)

type config struct {
	user     string `yaml:"user"`
	password string `yaml:"password"`
	dbname   string `yaml:"dbname"`
	host     string `yaml:"dbname"`
	port     int    `yaml:"port"`
	sslmode  string `yaml:"sslmode"`
	redisURL string `yaml:"redis_url"`
}

var schema = `
CREATE TABLE request (  
    id VARCHAR (50) NOT NULL PRIMARY KEY,
    driverid VARCHAR (50) NOT NULL ,
    riderid VARCHAR (50) NOT NULL,
    actual_price REAL,
    completed BOOLEAN NOT NULL,
    ride_time VARCHAR (50),
    distance REAL,
    average_time VARCHAR (50), 
    create_date VARCHAR (50)
);

CREATE TABLE place (
    id VARCHAR NOT NULL PRIMARY KEY ,
    origin_name VARCHAR (50),
    origin_latitude REAL,
    origin_longitude REAL,
    destination_name VARCHAR (50),
    destination_latitude REAL,
    destination_longitude REAL,
    
    FOREIGN KEY (id) REFERENCES request (id)
); 

CREATE TABLE ratings (
    rideid VARCHAR (50) NOT NULL,
    riderrating REAL,
    driverrating REAL
);
`

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
	createSchema()
}

func createSchema() {
	_, err := db.Exec(schema)
	if err != nil {
		log.Println(err)
	}

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


// DefaultService the CRUD
type DefaultService interface {
	Create() error
	Get() error
	Update() error
}

// Exist checks if the type exists
func Exist(u DefaultService) bool {
	err := u.Get()
	if err != nil {
		return false
	}
	return true
}
