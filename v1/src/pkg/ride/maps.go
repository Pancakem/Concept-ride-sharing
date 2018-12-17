package ride

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pancakem/rides/v1/src/pkg/store"
)

// file allows access to the google maps api

type val struct {
	Time map[string]int `json:"duration"`
}

var conv = fmt.Sprint

// InRange between the rider and driver and check if in range by time 10 minutes
// the success of this qualifies a driver to get a request sent to them
func InRange(l *store.LatLng, dl *store.DriverLocation, dur int) error {
	client := http.Client{Timeout: time.Second * 5}
	key := ""
	mapURL := "http://maps.googleapis.com/maps/api/directions/json?origin=" + conv(l.Lat) + "," + conv(l.Lng) + "&destination=" + conv(dl.Location.Lat) + "," + conv(dl.Location.Lng) + "&sensor=false&units=metric&mode=driving&key=" + key

	req, err := http.NewRequest("GET", mapURL, nil)
	if err != nil {
		// by returning err we move to the next driver
		// this sacrifice might save us time and we can get the next driver
		// the driver maybe placed in failed job queue for a next trial
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	values := val{}

	err = json.NewDecoder(resp.Body).Decode(&values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	tInt := values.Time["value"]
	duration := tInt
	dtime := duration / 60
	if dtime > dur {
		var er error
		return er
	}
	return nil
}


func ETA(origin *store.LatLng, destination *store.LatLng) float64 {
	client := http.Client{Timeout: time.Second * 5}
	key := ""
	mapURL := "http://maps.googleapis.com/maps/api/directions/json?origin=" + conv(origin.Lat) + "," + conv(origin.Lng) + "&destination=" + conv(destination.Lat) + "," + conv(destination.Lng) + "&sensor=false&units=metric&mode=driving&key=" + key

	req, err := http.NewRequest("GET", mapURL, nil)
	if err != nil {
		// by returning err we move to the next driver
		// this sacrifice might save us time and we can get the next driver
		// the driver maybe placed in failed job queue for a next trial
		return 0.0
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0.0
	}
	values := val{}

	err = json.NewDecoder(resp.Body).Decode(&values)
	if err != nil {
		return 0.0
	}
	defer resp.Body.Close()
	tInt := values.Time["value"]
	duration := tInt
	return float64(duration / 60)
}
