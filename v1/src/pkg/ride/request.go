package ride

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pancakem/rides/v1/src/pkg/store"
)

// ThisRequest is the request to be used
var (
	ThisRequest *store.DriverRequest
	upgrader    = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

// Match will pair a driver and a rider
func Match(hub *Hub, w http.ResponseWriter, r *http.Request) {
	var rr *store.RideRequest
	rid := make(chan *store.MatchResponse)
	rider, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Error making request", 500)
	}

	// read the request data
	err = rider.ReadJSON(rr)
	if err != nil {
		http.Error(w, "Couldn't parse data", 406)
	}

	// the ride consists of rider details sent to the driver
	ThisRequest := NewDriverRequest(rr)

	// use their id to get corresponding writers from the store

	distance := 5.0
	// get eight drivers that are in range of 5, 10 15
	// first try looks for 5
	for i := 0; i < 3; i++ {
		cli := store.GetRedisClient()
		dls := cli.SearchDrivers(8, rr.Origin.Lat, rr.Origin.Lng, distance)
		// send those drivers the request
		for _, val := range dls {
			conn := hub.Check(val.Name)
			go conn.Send()
			if conn == nil {
				continue
			}
			SendRideRequest(ThisRequest, conn)
			go conn.Read(rid)
			timer := time.NewTimer(time.Second * 15)
			select {
			case <-timer.C:
				continue
			case accepted := <-rid:
				// write to the rider the acceptance of their request
				data, _ := json.Marshal(accepted)
				conn.send <- data
				break
			}

		}
		distance += 5
	}

}

// spawn a go routine that will always listen for connection from the user and also
// another to write current location to the to the rider as the driver approaches
// when drive starts the
