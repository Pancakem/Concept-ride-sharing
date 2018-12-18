package ride

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/pancakem/rides/v1/src/pkg/store"
)

// Hub stores active drivers by their connections
type Hub struct {
	// connected clients (mostly drivers)
	clients map[string]*Client
	// add a new driver to the list
	register chan *aggreg
	// remove an inactive driver
	unregister chan *aggreg
}

type aggreg struct {
	id string
	*Client
}

// NewHub returns an instance of Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *aggreg),
		unregister: make(chan *aggreg),
	}
}

// Check finds a writer in the list of riders made online
func (h *Hub) Check(driverid string) *Client {
	if val, ok := h.clients[driverid]; ok {
		return val
	}
	return nil
}

// Run adds connections to hub and also removes them
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.id] = client.Client
		case client := <-h.unregister:
			if _, ok := h.clients[client.id]; ok {
				delete(h.clients, client.id)
				close(client.send)
			}
		}
	}
}

func (h *Hub) Read() {
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
							id:     key,
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
				}

			}

		}
	}

}
