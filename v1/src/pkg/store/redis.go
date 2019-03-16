package store

import (
	"sync"

	"github.com/go-redis/redis"
	"github.com/pancakem/rides/v1/src/pkg/common"
)

var once sync.Once
var redisClient *RedisClient

// RedisClient contains the pointer to the redis client
type RedisClient struct {
	*redis.Client
}

// GetRedisClient we get one instance of this in the application
func GetRedisClient() *RedisClient {
	conf := getConfig()
	once.Do(func() {
		client := redis.NewClient(&redis.Options{
			Addr:     conf.redisURL,
			Password: "",
			DB:       0, // using the default database
		})
		redisClient = &RedisClient{client}
	})
	_, err := redisClient.Ping().Result()
	if err != nil {
		common.Log.Printf("Could not connect to redis %v", err)
	}
	return redisClient
}

// AddDriverLocation adds driver id and LatLng to redis db
func (c *RedisClient) AddDriverLocation(dl *DriverLocation) {
	c.GeoAdd(
		dl.Vehicle,
		&redis.GeoLocation{Longitude: dl.Location.Lng, Latitude: dl.Location.Lat, Name: dl.DriverID},
	)
}

// RemoveDriverLocation from cache
func (c *RedisClient) RemoveDriverLocation(vehicletype, driverid string) {
	c.ZRem(vehicletype, driverid)
}

// SearchDrivers within a given radius
func (c *RedisClient) SearchDrivers(vehicletype string, limit int, lat, lng, r float64) []redis.GeoLocation {
	res, _ := c.GeoRadius(vehicletype, lng, lat, &redis.GeoRadiusQuery{
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
