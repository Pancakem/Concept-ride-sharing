package model

// BlackListed tokens
type BlackListed struct {
	Token string `bson:"token"`
}

func (bl *BlackListed) Create() error {
	_, err := sqldb.Exec("INSERT * INTO blacklisted VALUES(?)", bl.Token)
	return err
}

func (bl *BlackListed) Get() error {
	row := sqldb.QueryRow("SELECT * FROM blacklisted where token=?", bl.Token)
	var token string
	return row.Scan(&token)
}
