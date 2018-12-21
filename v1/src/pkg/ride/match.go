package ride

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pancakem/rides/v1/src/pkg/store"
)

// ThisRequest is the request to be used
var (
	upgrader = websocket.Upgrader{
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
	send := make(chan map[string]interface{})
	var rr store.RideRequest
	riderdata := make(chan map[string]interface{})
	rid := make(chan *store.MatchResponse)
	rider, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		http.Error(w, "Error making request", 500)
	}
	// read from rider
	go func() {
		defer func() {
			recover()
		}()
		for {
			ma := make(map[string]interface{})
			err := rider.ReadJSON(ma)
			if err != nil {
				http.Error(w, "Couldn't parse data", 406)
				continue
			}
			fmt.Println(ma)
			riderdata <- ma
		}
	}()

	go func() {
		for {
			select {
			case data := <-send:
				switch data["type"].(string) {
				case "finished":
					rider.WriteJSON(data)
				}

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

	cli := store.GetRedisClient()
	for i := 0; i < 3; i++ {

		dls := cli.SearchDrivers(8, rr.Origin.Lat, rr.Origin.Lng, distance)
		c := make(chan bool) // record accepted to exit outer loop

		// send those drivers the request
		for _, val := range dls {
			conn := hub.Check(val.Name)
			if conn == nil {
				continue
			}
			if conn.busy == false {
				go conn.Send()
				go conn.Read(rid, send)

				SendRideRequest(ThisRequest, rr.RiderID, conn)

				timer := time.NewTimer(time.Second * 15)
				select {
				case <-timer.C:
					continue
				case accepted := <-rid:
					// write to the rider the acceptance of their request
					accepted.ETA = ETA(&rr.Origin, &store.LatLng{Lat: accepted.Lat, Lng: accepted.Lng})
					rider.WriteJSON(accepted)
					conn.busy = true
					c <- true
					// write  to database
					store.Create(ThisRequest, val.Name, rr.RiderID)
					break
				}

			}
		}
		select {
		case <-c:
			break
		default:
			distance += 5
			continue
		}

	}

}
