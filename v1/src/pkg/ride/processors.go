package ride

import (
	"github.com/pancakem/rides/v1/src/pkg/store"
)

// saves the driver locations to the database
func saveToDB(dl *store.DriverLocation) {
	cli := store.GetRedisClient()
	cli.AddDriverLocation(dl)
}
