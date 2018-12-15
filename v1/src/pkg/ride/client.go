package ride

import (
	"encoding/json"

	"github.com/pancakem/rides/v1/src/pkg/store"
	"github.com/pancakem/swoop-rides-service/v1/src/pkg/model"

	"github.com/gorilla/websocket"
)

var maxMessageSize int64 = 512

// Client wraps the user connection and data to be sent through
type Client struct {
	conn *websocket.Conn
	send chan []byte
}

// NewClient returns a client instance
func NewClient(c *websocket.Conn) *Client {
	return &Client{
		conn: c,
	}
}

func (h *Hub) Read(rid chan *store.MatchResponse, an chan []byte) {

	for _, c := range h.clients {
		c.conn.SetReadLimit(maxMessageSize)

		for {
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					// unregister the client from hub
				}
				break
			}
			func(message []byte) {

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
					d.GetByID()

					rid <- &store.MatchResponse{
						LatLng: acc.Location,
						Driver: *d,
					}
				case "cancelled":
					// cancelled should contain ride id
					driverCancel(ma["id"].(string), ma["time"].(float64), ma["distance"].(float64), an)
					c.send <- <-an

				}
			}(message)

		}
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
