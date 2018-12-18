package ride

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/pancakem/rides/v1/src/pkg/store"
	"github.com/pancakem/user-service/v1/src/pkg/model"
)

var maxMessageSize int64 = 512

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

func (c *Client) Read(rid chan *store.MatchResponse, an chan map[string]interface{}) {
	for {
		c.conn.SetReadLimit(maxMessageSize)

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
				d := &model.Driver{ID: acc.DriverID}
				err := d.GetByID()
				if err != nil {
					log.Println((err))
				}
				err = d.Vehicle.Get()
				if err != nil {
					log.Println((err))
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
				finished("cancelled", ma["id"].(string), ma["time"].(float64), ma["distance"].(float64), an, c, false)
				c.busy = false
				break
				return

			case "finished":
				finished("finished", ma["id"].(string), ma["time"].(float64), ma["distance"].(float64), an, c, true)
				c.busy = false
				break
				return

			case "rating":
				store.AddRiderRating(ma["id"].(string), ma["rating"].(float32))

			}

		}
	}

}
