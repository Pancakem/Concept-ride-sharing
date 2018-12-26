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

// Get checks for the token if exists in the database
func (bl *BlackListed) Get() bool {
	row := sqldb.QueryRow("SELECT * FROM blacklisted where token=?", bl.Token)
	var token string
	err := row.Scan(&token)
	if err == sql.ErrNoRows || err != nil {
		return false
	}
	return true
}
