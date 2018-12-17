package ride

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pancakem/rides/v1/src/pkg/store"
	"github.com/pancakem/swoop-rides-service/v1/src/pkg/model"
	"log"
)

var maxMessageSize int64 = 512

// Client wraps the user connection and data to be sent through
type Client struct {
	conn *websocket.Conn
	send chan []byte
	booked bool
}

// NewClient returns a client instance
func NewClient(c *websocket.Conn) *Client {
	return &Client{
		conn: c,
	}
}

func (h *Hub) Read(rid chan *store.MatchResponse, an chan []byte) {
	fmt.Println("started h.Read")
	cli := store.GetRedisClient()
	if len(h.clients) > 0 {
		for key, c := range h.clients {
			c.conn.SetReadLimit(maxMessageSize)

			for {
				_, message, err := c.conn.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						// unregister the client from hub
						h.unregister <- &aggreg{
							id: key,
							Client: c,
						}
						// delete location from redis
						cli.RemoveDriverLocation(key)

					}
					break
				}

				ma := make(map[string]interface{})
				json.Unmarshal(message, &ma)

				switch ma["type"] {
				case "driverlocation":
					dl := &store.DriverLocation{}
					json.Unmarshal(message, dl)
					saveToDB(dl)
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
						Type: "accepted",
						LatLng: acc.Location,
						Name: d.FullName,
						PhoneNumber:d.Phonenumber,
						ImageURL: d.ProfileImage,
						VehicleColor:d.Vehicle.Color,
						VehicleModel:d.Vehicle.Model,
						VehiclePlate:d.Vehicle.PlateNumber,
					}
				case "cancelled":
					// cancelled should contain ride id
					finished("cancelled", ma["id"].(string), ma["time"].(float64), ma["distance"].(float64), an, false)
					c.booked = false
					c.send <- <-an
				case "finished":
					finished("finished", ma["id"].(string), ma["time"].(float64), ma["distance"].(float64), an, true)
					c.booked = false
					c.send <- <-an

				case "rating":
					store.AddRiderRating(ma["id"].(string), ma["rating"].(float32))

				}

			}
		}

	}



}

// Send any data written to send channel to driver
// for low priority  data
func (c *Client) Send() {
	fmt.Println("Started c.Send")
	for {
		select {
		case data := <-c.send:
			c.conn.WriteJSON(data)
		}
	}

}
