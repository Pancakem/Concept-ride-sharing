package model

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" //
	"github.com/pancakem/rides/v1/src/pkg/setting"
)

var schema = ``

var sqldb *sql.DB

func setSQLEngine() {
	conf := setting.Config{}
	conf.Get()
	dataSourceName := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		conf.Host, conf.Port, conf.User, conf.Password, conf.DatabaseName, conf.SSLMode)

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
}

// Exist checks if the type exists
func Exist(u DefaultService) bool {
	err := u.Get()
	if err != nil {
		return false
	}
	return true
}
