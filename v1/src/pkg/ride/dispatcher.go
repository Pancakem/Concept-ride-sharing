package ride

import (
	"encoding/json"

	"github.com/pancakem/rides/v1/src/pkg/store"
	"github.com/pancakem/swoop-rides-service/v1/src/pkg/model"
)

// dispatch finds out what kind of event the ride has requested for
//it then calls the appropriate function to handle it

func dispatch(message []byte, rid chan *store.MatchResponse) {

	ma := make(map[string]interface{})
	json.Unmarshal(message, &ma)

	switch ma["type"] {
	case "driverlocation":
		dl := &store.DriverLocation{}
		json.Unmarshal(message, dl)
		saveToDB(dl)
	case "accepted":
		// when a ride is accepted
		// the driver details are found out
		// and his/her location
		// this data is sent to the rider channel which
		// will trigger the send response to them
		acc := &store.Accepted{}
		json.Unmarshal(message, acc)
		d := &model.Driver{ID: acc.DriverID}
		d.GetByID()

		rid <- &store.MatchResponse{
			LatLng: acc.Location,
			Driver: *d,
		}
	
		

	}
}
