package store

import (
	"time"
)

// Request is the finalized data about a completed or not ride
type request struct {
	ID          string     `db:"id"`
	CreateDate  *time.Time `db:"create_time"`
	DriverID    string     `db:"driver"`
	RiderID     string     `db:"rider"`
	Origin      string     `db:"origin"`
	Destination string     `db:"destination"`
	ActualPrice float32    `db:"actualprice"`
	Completed   bool       `db:"completed"`
	Ratings     float32    `db:"rating"`
	AverageTime *time.Time `db:"average_time_of_booking"`
	Canceled    bool       `db:"canceled"`
}

// Create commits a ride request into the database
func Create(reques *RideRequest, requestid string) error {
	// decode auth key for rider id
	// use googlemaps api to get current location name and destination name
	r := &request{ID: requestid}
	return create(r)
}
func create(r *request) error {
	stmt, err := db.Prepare("INSERT * INTO request")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(r)
	if err != nil {
		return err
	}
	return nil
}

func read(r *request) error {
	row := db.QueryRow("SELECT * FROM request where id=$1", r.ID)
	return row.Scan(&r)
}

// repair
func update(r *request) error {
	stmt, err := db.Prepare("UPDATE request SET * WHERE id=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(&r)
	return err
}

func updateRating(reqID string, rating float32) error {
	_, err := db.Exec("UPDATE request SET rating=$1 where id=$2", rating, reqID)
	return err
}
