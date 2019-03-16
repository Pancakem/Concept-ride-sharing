package store

import "github.com/pancakem/rides/v1/src/pkg/common"

// Rider definition
type Rider struct {
	ID            string `db:"id" json:"_id"`
	FullName      string `db:"fullname" json:"fullname"`
	Email         string `db:"email" json:"email"`
	Phonenumber   string `db:"phonenumber" json:"phonenumber"`
	Password      string `db:"password" json:"password"`
	PaymentMethod string `db:"paymentmethod" json:"paymentmethod"`
	IsActive      bool   `db:"is_active"`
}

// GetByID retreives rider by id
func (u *Rider) GetByID() error {
	row := db.QueryRow("SELECT * FROM rider WHERE id=$1", u.ID)
	err := row.Scan(&u.ID, &u.FullName, &u.Email,
		&u.Phonenumber, u.Password, u.IsActive, u.PaymentMethod)
	return err
}

// Get is used by loggin retrieves by email or phonenumber
func (u *Rider) Get() error {
	if u.Email != "" {
		row := db.QueryRow("SELECT id, fullname, password_ FROM rider WHERE email=$1", u.Email)
		err := row.Scan(&u.ID, &u.FullName, &u.Password)
		return err
	}
	row := db.QueryRow("SELECT id, password_ FROM rider WHERE phonenumber=$1", u.Phonenumber)
	err := row.Scan(&u.ID, &u.Password)
	return err
}

// Create inserts rider into db
func (u *Rider) Create() error {
	stmt, err := db.Prepare(`INSERT INTO rider (id, fullname, email, 
		phonenumber, password_, paymentmethod, isactive) VALUES($1,$2,$3,$4,$5,$6,$7)`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(u.ID, u.FullName, u.Email, u.Phonenumber, u.Password, u.PaymentMethod, u.IsActive)
	return err
}

// Update up
func (u *Rider) Update() error {
	stmt, err := db.Prepare("UPDATE rider SET email=$1 phonenumber=$2 payment=$3 WHERE id=$4")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(u.Email, u.Phonenumber, u.PaymentMethod, u.ID)
	return err
}

// UpdatePassword when a user wants
func (u *Rider) UpdatePassword() error {
	_, err := db.Exec("UPDATE rider SET password=$1 WHERE id=$2", u.Password, u.ID)
	return err
}

// Delete a rider from records
func (u *Rider) Delete() error {
	_, err := db.Exec("DELETE * FROM rider WHERE id=$1", u.ID)
	return err
}

// GetAllRider fetches all
func GetAllRider() []Rider {
	riders := []Rider{}
	rows, err := db.Query("SELECT * FROM rider")
	if err != nil {
		common.Log.Println(err)
	}
	for rows.Next() {
		r := Rider{}
		if err = rows.Scan(&r.ID, &r.FullName, &r.Email,
			&r.Phonenumber, r.Password, r.IsActive, r.PaymentMethod); err != nil {
			common.Log.Println(err)
		}
		riders = append(riders, r)

	}
	return riders
}

// LoginForm used to abstract loggin
type LoginForm struct {
	Email       string
	Phonenumber string
	Password    string
}

// WhatIsNil Find which parameter the user wants to loggin with
func (lf *LoginForm) WhatIsNil() (string, string) {
	if lf.Phonenumber != "" {
		return "phonenumber", lf.Phonenumber
	}
	return "email", lf.Email
}

// Get abstracts the loggin to either rider or driver
func (lf *LoginForm) Get() (interface{}, error) {
	ty, val := lf.WhatIsNil()
	rid := Rider{}
	driv := Driver{}
	if ty == "phonenumber" {
		rid.Phonenumber = val
		err := rid.Get()
		if err != nil {
			common.Log.Println(err)
			driv.Phonenumber = val
			err = driv.Get()
			if err != nil {
				return nil, err
			}
			return driv, nil
		}
	}
	rid.Email = val
	err := rid.Get()
	if err != nil {
		common.Log.Println(err)
		driv.Email = val
		err = driv.Get()
		if err != nil {
			return nil, err
		}
		return driv, nil
	}

	return rid, nil
}
