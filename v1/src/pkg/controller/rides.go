package controller

import (
	"net/http"

	"github.com/pancakem/rides/v1/src/pkg/ride"
)

var hub *ride.Hub

// BookRide allows the rider to request a ride
func BookRide(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ride.Match(hub, w, r)
}

// GetLocation takes locations of drivers
func GetLocation(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	hub = ride.NewHub()
	ride.GetLocation(hub, w, r)
}
