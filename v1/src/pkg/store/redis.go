package store

import (
	"log"
	"sync"

	"github.com/go-redis/redis"
)

const (
	key = "drivers"
)

var once sync.Once
var redisClient *RedisClient

// RedisClient contains the pointer to the redis client
type RedisClient struct {
	*redis.Client
}

// GetRedisClient we get one instance of this in the application
func GetRedisClient() *RedisClient {
	once.Do(func() {
		client := redis.NewClient(&redis.Options{
			Addr:     "0.0.0.0:32768",
			Password: "",
			DB:       0, // using the default database
		})
		redisClient = &RedisClient{client}
	})
	_, err := redisClient.Ping().Result()
	if err != nil {
		log.Fatalf("Could not connect to redis %v", err)
	}
	return redisClient
}

// AddDriverLocation adds driver id and LatLng to redis db
func (c *RedisClient) AddDriverLocation(dl *DriverLocation) {
	c.GeoAdd(
		key,
		&redis.GeoLocation{Longitude: dl.Location.Lng, Latitude: dl.Location.Lat, Name: dl.DriverID},
	)
}

// RemoveDriverLocation from cache
func (c *RedisClient) RemoveDriverLocation(driverid string) {
	c.ZRem(key, driverid)
}

// SearchDrivers within a given radius
func (c *RedisClient) SearchDrivers(limit int, lat, lng, r float64) []redis.GeoLocation {
	res, _ := c.GeoRadius(key, lng, lat, &redis.GeoRadiusQuery{
		Radius:      r,
		Unit:        "km",
		WithCoord:   true,
		WithDist:    true,
		WithGeoHash: true,
		Sort:        "ASC",
		Count:       limit,
	}).Result()
	return res
}

// UpdateLocation updates the driver location
func (c *RedisClient) UpdateLocation(lng, lat float64, driverid string) {

}

// Exist confirms if a driver exists if not a new entry is made else
// the location is updated
func (c *RedisClient) Exist(driverid string) bool {
	return true
}
