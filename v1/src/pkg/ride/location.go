package ride

import (
	"net/http"

	"github.com/pancakem/rides/v1/src/pkg/store"
)

// GetLocation of the drivers by creating websockets
func GetLocation(hub *Hub, w http.ResponseWriter, r *http.Request) {
	driver, err := upgrader.Upgrade(w, r, nil) // later update response headers
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
	}

	client := &Client{conn: driver, send: make(chan []byte, 256)}
	dl := &store.DriverLocation{}

	client.conn.ReadJSON(dl)

	go hub.run()

	rediscli := store.GetRedisClient()
	rediscli.AddDriverLocation(dl)
	a := aggreg{id: dl.DriverID, Client: client}
	hub.register <- &a
}
