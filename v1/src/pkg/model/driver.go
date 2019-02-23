package model

import (
	"github.com/pancakem/rides/v1/src/pkg/common"
)

// Driver definition
type Driver struct {
	ID           string `json:"_id"`
	FullName     string `json:"fullname"`
	Email        string `djson:"email"`
	Phonenumber  string `json:"phonenumber"`
	NationalID   string `json:"national_id"`
	LicenseNo    string `json:"license_no"`
	ProfileImage string `json:"profile_image"`
	Vehicle
	IsActive bool
	Password string `json:"password"`
}

// Vehicle definition
type Vehicle struct {
	ID          string
	DriverID    string
	Owner       string
	ImageURL    string
	Color       string
	Model       string
	PlateNumber string
	TypeOf      string
}

func GetDriverByID(id string) (*Driver, error) {
	row := sqldb.QueryRow(`SELECT id, fullname, email, phonenumber, password_, isactive, 
	national_id, license_no, profile_image,vehicle_id FROM driver WHERE id=$1`, id)

	u := &Driver{}

	err := row.Scan(&u.ID, &u.FullName, &u.Email, &u.Phonenumber, &u.Password,
		&u.IsActive, &u.NationalID, &u.LicenseNo, &u.ProfileImage, &u.Vehicle.ID)

	if err != nil {
		return nil, err
	}
	return u, err
}

func (u *Driver) Get() error {
	if u.Email != "" {
		row := sqldb.QueryRow("SELECT id, fullname, password_ FROM driver WHERE email=$1", u.Email)
		err := row.Scan(&u.ID, &u.FullName, &u.Password)
		return err
	}
	row := sqldb.QueryRow("SELECT id, fullname, password_ FROM driver WHERE phonenumber=$1", u.Phonenumber)
	err := row.Scan(&u.ID, &u.FullName, &u.Password)
	return err
}

func (u *Driver) Create() error {
	stmt, err := sqldb.Prepare(`INSERT INTO driver (id, fullname, email, 
		phonenumber, password_, profile_image, license_no, national_id, isactive)VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)`)
	if err != nil {
		return err
	}
	u.IsActive = true
	_, err = stmt.Exec(u.ID, u.FullName, u.Email, u.Phonenumber, u.Password, u.ProfileImage, u.LicenseNo, u.NationalID, u.IsActive)
	return err
}

func (u *Driver) Update() error {
	stmt, err := sqldb.Prepare("UPDATE rider SET email=$1 phonenumber=$2 payment=$3 WHERE id=$4")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(u.Email, u.Phonenumber, u.Vehicle.ID, u.ID)
	return err

}

func DeleteVehicle(id string) error {
	_, err := sqldb.Exec("DELETE FROM driver WHERE id=$1", id)
	return err
}

func GetAllDriver() []Driver {
	drivers := []Driver{}
	rows, err := sqldb.Query("SELECT * FROM driver")
	if err != nil {
		common.Log.Println(err)
	}
	for rows.Next() {
		d := Driver{}
		if err = rows.Scan(d.ID, d.FullName, d.Email, d.Phonenumber, d.Password, d.ProfileImage, d.LicenseNo, d.NationalID, d.IsActive); err != nil {
			common.Log.Println(err)
		}
		drivers = append(drivers, d)

	}

	return drivers
}


// Lock bans driver from platform
func (u *Driver) Lock() {
	u.IsActive = false
	u.Update()
}

// Create adds a vehicle row to the database
func (v *Vehicle) Create() error {
	stmt, err := sqldb.Prepare(`INSERT INTO vehicles (id, driver_id, owner, 
		typeof, color, model, plate_no, image_url) VALUES($1,$2,$3,$4,$5,$6,$7,$8)`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(v.ID, v.DriverID, v.Owner, v.TypeOf, v.Color, v.Model, v.PlateNumber, v.ImageURL)
	return err

}

// GetVehicle returns a pointer with data from database row
func GetVehicle(id string) (*Vehicle, error) {
	var v *Vehicle
	row := sqldb.QueryRow("SELECT color, model, plate_no FROM vehicles WHERE id=$1", id)

	if err := row.Scan(v.Color, v.Model, v.PlateNumber); err != nil {
		return v, err
	}
	return v, nil
}

// GetVehicleType returns a string that is the type of vehicle a driver has
func GetVehicleType(id string) string {
	var typevehicle string

	row := sqldb.QueryRow("SELECT typeof FROM vehicles WHERE driver_id=$1", id)

	if err := row.Scan(typevehicle); err != nil {
		return ""
	}
	return typevehicle
}
