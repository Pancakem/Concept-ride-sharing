package store

import "github.com/pancakem/rides/v1/src/pkg/common"

// AddDriverRating adds rating from rider
func AddDriverRating(rideid string, rating float32) error {

	sql := `INSERT INTO rating rideid, driverrating
	VALUES($1,$2)`

	stmt, err := db.Prepare(sql)
	if err != nil {
		common.Log.Println(err)
	}
	_, err = stmt.Exec(rideid, rating)
	return err
}

// AddRiderRating adds rating from driver
func AddRiderRating(rideid string, rating float32) error {

	sql := `INSERT INTO rating rideid, riderrating
	VALUES($1,$2)`

	stmt, err := db.Prepare(sql)
	if err != nil {
		common.Log.Println(err)
	}
	_, err = stmt.Exec(rideid, rating)
	return err
}
