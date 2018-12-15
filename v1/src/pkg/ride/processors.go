package ride

import (
	"encoding/json"

	"github.com/pancakem/rides/v1/src/pkg/payment"
	"github.com/pancakem/rides/v1/src/pkg/store"
)

// saves the driver locations to the database
func saveToDB(dl *store.DriverLocation) {
	cli := store.GetRedisClient()
	cli.AddDriverLocation(dl)
}

// when ride is cancelled receive a canceled and ride stats time and distance float64
func driverCancel(rideid string, time float64, distance float64, an chan []byte) {
	// calculate the ride amount
	price := payment.Price.Calculate(time/60, distance)
	ma := make(map[string]interface{})
	ma["type"] = "cancelled"
	ma["price"] = price

	// write to channel cancelled and trigger a send json to rider writer
	data, _ := json.Marshal(ma)
	an <- data

	defer func() {
		// update the database table
		store.Update(rideid, time, price, false)
	}()
}

func riderCancel(rideid string, time float64, distance float64, hub *Hub, an chan []byte) {
	// calculate the ride amount
	price := payment.Price.Calculate(time/60, distance)
	ma := make(map[string]interface{})
	ma["type"] = "cancelled"
	ma["price"] = price

	// read db and find the connection
	r := store.Read(rideid)

	c := hub.Check(r.DriverID)
	c.conn.WriteJSON(ma)
	data, _ := json.Marshal(ma)
	an <- data
	defer func() {
		store.Update(rideid, time, price, false)
	}()

}
