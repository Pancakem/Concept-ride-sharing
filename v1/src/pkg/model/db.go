package model

import (
	"database/sql"
	"fmt"
	"log"

	// _ "github.com/mattn/go-sqlite3"
	_ "github.com/lib/pq" // uncomment for production environment
	"github.com/pancakem/rides/v1/src/pkg/setting"
)

var sqldb *sql.DB

func setSQLEngine() {
	conf := setting.Config{}
	conf.Get()
	dataSourceName := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		conf.Host, conf.Port, conf.User, conf.Password, conf.DatabaseName, conf.SSLMode)
	// toDb, err := sqlx.Open("postgres", dataSourceName)
	toDb, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Println("Failed to connect to database", err)
	}
	sqldb = toDb
	_, err = sqldb.Exec(schema)
	if err != nil {
		log.Println(err)
	}

}

func init() {
	setSQLEngine()
}

// DefaultService the CRUD
type DefaultService interface {
	Create() error
	Get() error
	Update() error
	Delete() error
	GetAll() interface{}
}

var schema = `
CREATE TABLE rider (
    id VARCHAR NOT NULL PRIMARY KEY,
    fullname VARCHAR NOT NULL,
    email VARCHAR,
    phonenumber VARCHAR,
    password_ VARCHAR,
    isactive BOOLEAN,
    paymentmethod VARCHAR,
    create_date VARCHAR
);


CREATE TABLE driver (
    id VARCHAR NOT NULL PRIMARY KEY,
    fullname VARCHAR NOT NULL,
    email VARCHAR ,
    phonenumber VARCHAR ,
    password_ VARCHAR,
    isactive BOOLEAN,
    national_id VARCHAR,
    license_no VARCHAR ,
    profile_image VARCHAR,
    vehicle_id VARCHAR,
    create_date VARCHAR
);


CREATE TABLE vehicles (
    id VARCHAR NOT NULL PRIMARY KEY,
    driver_id VARCHAR,
    image_url VARCHAR,
    owner_ VARCHAR,
    color VARCHAR,
    typeof VARCHAR,
    plate_no VARCHAR,
    model VARCHAR
);

-- expired jwt tokens
CREATE TABLE blacklisted (
    token VARCHAR NOT NULL
);
`
