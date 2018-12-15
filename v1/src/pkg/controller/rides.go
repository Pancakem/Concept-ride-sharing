package controller

import (
	"fmt"
	"net/http"

	"github.com/pancakem/rides/v1/src/pkg/ride"
)

var hub *ride.Hub

// BookRide allows the rider to request a ride
func BookRide(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	fmt.Println("in match")
	ride.Match(hub, w, r)
}

// GetLocation takes locations of drivers
func GetLocation(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	hub = ride.NewHub()
	fmt.Println("in loco")
	ride.GetLocation(hub, w, r)
}
