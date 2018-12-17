package store

import "log"

// AddDriverRating adds rating from rider
func AddDriverRating(rideid string, rating float32) error {

	sql := `INSERT INTO rating rideid, driverrating
	VALUES($1,$2)`

	stmt, err := db.Prepare(sql)
	if err != nil {
		log.Println(err)
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
		log.Println(err)
	}
	_, err = stmt.Exec(rideid, rating)
	return err
}
