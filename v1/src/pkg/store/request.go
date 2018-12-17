package store

import (
	"time"
)

// Request is the finalized data about a completed or not ride
type Request struct {
	ID          string     `db:"id"`
	DriverID    string     `db:"driver"`
	RiderID     string     `db:"rider"`
	Origin      LatLng     `db:"origin"`
	Destination LatLng     `db:"destination"`
	ActualPrice float32    `db:"actual_price"`
	Completed   bool       `db:"completed"`
	RideTime    float64    `db:"ride_time"` // ride time in seconds
	CreateDate  *time.Time `db:"create_time"`
	AverageTime *time.Time `db:"average_time_of_booking"`
}

// Create commits a ride Request into the database
func Create(reques *DriverRequest, driverid, riderid string) error {
	req := &Request{
		ID:          reques.RequestID,
		DriverID:    driverid,
		RiderID:     riderid,
		Origin:      reques.Origin,
		Destination: reques.Destination,
		Completed:   false,
	}
	err := createPlace(req)
	err = createRequest(req)
	return err
}

// writes to Request table
func createRequest(r *Request) error {
	// insert into Request table
	sql := `INSERT INTO request
	id, 
	driverid,
	riderid,
	completed,
  	VALUES($1, $2, $3, $4)
	`

	stmt, err := db.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(r.ID, r.DriverID, r.RiderID, r.Completed)
	if err != nil {
		return err
	}
	return nil
}

// writes to place table
func createPlace(r *Request) error {
	sql := `INSERT INTO place
	id, 
	origin_name,
	origin_latitude,
	origin_longitude,
	destination_name,
	destination_latitude,
	destination_longitude, 
	`
	stmt, err := db.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(r.ID, r.Origin.Lat, r.Origin.Lng, r.Origin.PlaceName,
		r.Destination.Lat, r.Destination.Lng, r.Destination.PlaceName)
	if err != nil {
		return err
	}
	return nil
}
func read(id string) *Request {
	r := &Request{}

	query := `SELECT 
	request.driver, 
	request.rider, 
	request.completed, 
	request.actual_price, 
	request.create_time, 
	request.average_time,
	place.origin_name,
	place.origin_latitude,
	place.origin_longitude,
	place.destination_name,
	place.destination_latitude,
	place.destination_longitude
	FROM request
	FULL OUTER JOIN place ON request.id=place.id;
	`
	row := db.QueryRow(query, id)
	err := row.Scan(r.DriverID, r.RiderID, r.Completed,
		r.ActualPrice, r.CreateDate, r.AverageTime,
		r.Origin.PlaceName, r.Origin.Lat, r.Origin.Lng,
		r.Destination.PlaceName, r.Destination.Lat, r.Destination.Lng,
	)
	if err != nil {
		return nil
	}

	return r
}

// Read the database
func Read(id string) *Request {
	return read(id)
}

// repair
func update(id string, ridetime float64, price float64, status bool) error {
	sql := `UPDATE request SET 
	completed=$1
	actual_price=$2
	ride_time=$3
	WHERE id=$4`
	stmt, err := db.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(status, price, ridetime, id)
	return err
}

// Update the request when its done
func Update(id string, ridetime float64, price float64, status bool) error {
	return update(id, ridetime, price, status)
}
