package ride

import (
	"github.com/pancakem/rides/v1/src/pkg/payment"
	"github.com/pancakem/rides/v1/src/pkg/store"
)

// saves the driver locations to the database
func saveToDB(dl *store.DriverLocation) {
	cli := store.GetRedisClient()
	cli.AddDriverLocation(dl)
}

func riderCancel(rideid string, time float64, distance float64, hub *Hub, an chan map[string]interface{}) {
	// read db and find the connection
	r := store.Read(rideid)
	// calculate the ride amount
	pricing := payment.Vehicle(r.VehicleType)

	ma := make(map[string]interface{})
	ma["type"] = "cancelled"

	price := pricing.Calculate(time/60, distance)
	ma["price"] = price

	c := hub.Check(r.DriverID)
	c.conn.WriteJSON(ma)
	an <- ma
	defer func() {
		store.Update(rideid, time, price, false)
	}()

}

func finished(typeo, rideid string, time float64, distance float64, vehicleType string, an chan map[string]interface{}, c *Client, status bool) {
	// calculate the ride amount
	pricing := payment.Vehicle(vehicleType)
	price := pricing.Calculate(time/60, distance)
	ma := make(map[string]interface{})
	ma["type"] = typeo
	ma["price"] = price

	// write to channel cancelled and trigger a send json to rider writer
	an <- ma
	c.conn.WriteJSON(ma)
	defer func() {
		// update the database table
		store.Update(rideid, time, price, status)
	}()
}

func rateRider(riderID, rideID string, rating float32) {
	// write database
}
