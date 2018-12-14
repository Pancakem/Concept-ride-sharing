package ride

import (
	"github.com/pancakem/rides/v1/src/pkg/store"

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

func (c *Client) Read(rid chan *store.MatchResponse) {

	c.conn.SetReadLimit(maxMessageSize)

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// unregister the client from hub
			}
			break
		}
		dispatch(message, rid)

	}

}

func (c *Client) Send() {
	for {
		select {
		case data := <-c.send:
			c.conn.WriteJSON(data)
		}
	}

}
