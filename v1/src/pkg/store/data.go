package store

// LatLng represents a location onnthe surface of earth
type LatLng struct {
	Lat       float64 `json:"lat"`
	Lng       float64 `json:"lng"`
	PlaceName string  `json:"place"`
}

// RideRequest takes in the json
type RideRequest struct {
	// rider_id key
	RiderID string `json:"rider_id"`
	// origin coordinates
	Origin LatLng `json:"origin"`
	// destination coordinates
	Destination LatLng `json:"destination"`
}

// DriverRequest json serial is used to send driver the request
type DriverRequest struct {
	// rider location
	RequestID string `json:"id"`
	// origin coordinates
	Origin LatLng `json:"origin"`
	// destination coordinates
	Destination LatLng `json:"destination"`
}

// Accepted request
type Accepted struct {
	DriverID string `json:"driver_id"`
	Location LatLng `json:"location"`
}

// DriverLocation wrapper
type DriverLocation struct {
	DriverID string `json:"driver_id"`
	Location LatLng `json:"location"`
}

// MatchResponse wrapper
type MatchResponse struct {
	Type string `json:"type"`
	LatLng `json:"location"`
	Name string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	ImageURL string `json:"image_url"`
	VehicleModel string `json:"vehicle_model"`
	VehicleColor string `json:"vehicle_color"`
	VehiclePlate string `json:"vehicle_plate"`
	ETA float64 `json:"ETA"`
}
