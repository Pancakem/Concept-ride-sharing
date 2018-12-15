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
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// Match will pair a driver and a rider
func Match(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// when written to this channel send data to the rider
	// for low priority data
	send := make(chan []byte)
	var rr store.RideRequest
	riderdata := make(chan map[string]interface{})
	rid := make(chan *store.MatchResponse)
	rider, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		http.Error(w, "Error making request", 500)
	}
	// read from rider
	go func() {
		for {
			ma := make(map[string]interface{})
			err = rider.ReadJSON(ma)
			if err != nil {
				http.Error(w, "Couldn't parse data", 406)
				continue
			}
			riderdata <- ma
		}

	}()

	go func() {
		for {
			select {
			case data := <-send:
				rider.WriteMessage(2, data)
			case x := <-riderdata:
				switch x["type"].(string) {
				case "cancelled":
					riderCancel(x["id"].(string), x["time"].(float64), x["distance"].(float64), hub, send)
				case "request":
					dta, _ := json.Marshal(x)
					json.Unmarshal(dta, &rr)
				case "rating":
					store.AddDriverRating(x["id"].(string), x["rating"].(float32))
				}
			}
		}

	}()

	// the ride consists of rider details sent to the driver
	ThisRequest := NewDriverRequest(&rr)
	// the first range of distance to try
	distance := 5.0
	// get eight drivers that are in range of 5, 10 15
	// first try looks for 5
	go hub.Read(rid, send)
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

			timer := time.NewTimer(time.Second * 15)
			select {
			case <-timer.C:
				continue
			case accepted := <-rid:
				// write to the rider the acceptance of their request
				data, _ := json.Marshal(accepted)
				rider.WriteMessage(2, data)
				break
			}

		}
		distance += 5
	}

}

// spawn a go routine that will always listen for connection from the user and also
// another to write current location to the to the rider as the driver approaches
// when drive starts the
