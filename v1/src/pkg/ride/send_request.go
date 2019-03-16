package ride

import (
	"github.com/pancakem/rides/v1/src/pkg/store"

	"github.com/pancakem/rides/v1/src/pkg/common"
)

// NewDriverRequest creates a driver request for the job queue
func NewDriverRequest(r *store.RideRequest) *store.DriverRequest {
	id, err := common.NewID()
	if err != nil {
		common.Log.Println("Failed to get uuid", err)
		return nil
	}
	return &store.DriverRequest{
		RequestID:   id,
		Origin:      r.Origin,
		Destination: r.Destination,
	}
}

// SendRideRequest encodes and sends a driver
func SendRideRequest(dr *store.DriverRequest, riderid string, w *Client) {
	// create a rider instance to get data from db
	x := store.Rider{ID: riderid}
	x.GetByID()

	ma := make(map[string]interface{})
	ma["id"] = dr.RequestID
	ma["name"] = x.FullName
	ma["phone_number"] = x.Phonenumber
	ma["origin"] = dr.Origin
	ma["destination"] = dr.Origin
	ma["rating"] = 5.0 // add datase field on user to write rating

	w.conn.WriteJSON(ma)
}
