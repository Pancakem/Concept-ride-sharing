package ride

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/pancakem/rides/v1/src/pkg/common"
	"github.com/pancakem/rides/v1/src/pkg/store"
)

func maxMessageSize() int64 { return 512 }

// Client wraps the user connection and data to be sent through
type Client struct {
	conn *websocket.Conn
	send chan []byte
	busy bool
}

// NewClient returns a client instance
func NewClient(c *websocket.Conn) *Client {
	return &Client{
		conn: c,
	}
}

// Send any data written to send channel to driver
// for low priority  data
func (c *Client) Send() {
	for {
		select {
		case data := <-c.send:
			c.conn.WriteJSON(data)
		}
	}

}

// Read data from a client
func (c *Client) Read(rid chan *store.MatchResponse, an chan map[string]interface{}) {

	c.conn.SetReadLimit(maxMessageSize())

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			// do something clever
			break
		}

		ma := make(map[string]interface{})
		json.Unmarshal(message, &ma)

		switch ma["type"] {
		case "driverlocation":
			if c.busy {
				an <- ma
			} else {
				dl := &store.DriverLocation{}
				json.Unmarshal(message, dl)
				saveToDB(dl)
			}

		case "accepted":
			// when a ride is accepted
			// the driver details are found out
			// and his/her location
			// this data is sent to the rider channel which
			// will trigger the send response to them
			acc := &store.Accepted{}
			json.Unmarshal(message, acc)
			d, err := store.GetDriverByID(acc.DriverID)
			if err != nil {
				common.Log.Println((err))
			}
			pointerToVehicle, err := store.GetVehicle(d.ID)
			d.Vehicle = *pointerToVehicle
			if err != nil {
				common.Log.Println((err))
			}

			rid <- &store.MatchResponse{
				Type:         "accepted",
				LatLng:       acc.Location,
				Name:         d.FullName,
				PhoneNumber:  d.Phonenumber,
				ImageURL:     d.ProfileImage,
				VehicleColor: d.Vehicle.Color,
				VehicleModel: d.Vehicle.Model,
				VehiclePlate: d.Vehicle.PlateNumber,
			}
		case "cancelled":
			// cancelled should contain ride id
			finished("cancelled", ma["id"].(string), ma["time"].(float64), ma["distance"].(float64), ma["vehicle_type"].(string), an, c, false)
			c.busy = false
			return

		case "finished":
			finished("finished", ma["id"].(string), ma["time"].(float64), ma["distance"].(float64), ma["vehicle_type"].(string), an, c, true)
			c.busy = false
			return

		case "rating":
			store.AddRiderRating(ma["id"].(string), ma["rating"].(float32))

		}

	}
}
