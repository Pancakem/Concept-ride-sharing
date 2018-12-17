package ride

import (
	"fmt"
	"github.com/gorilla/websocket"
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
