package model

import (
	"database/sql"
)

// BlackListed tokens
type BlackListed struct {
	Token string `bson:"token"`
}

// Create adds a used or expired token to db so that it can't be misused
func (bl *BlackListed) Create() error {
	_, err := sqldb.Exec("INSERT * INTO blacklisted VALUES(?)", bl.Token)
	return err
}

// GetToken checks for the token if exists in the database
func GetToken(token string) bool {
	row := sqldb.QueryRow("SELECT * FROM blacklisted where token=?", token)
	var toke string
	err := row.Scan(&toke)
	if err == sql.ErrNoRows || err != nil {
		return false
	}
	return true
}
