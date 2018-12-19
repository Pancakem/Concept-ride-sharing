package common

import (
	"log"
	"os"
)

var (
	Log *log.Logger
)

func init() {
	// c := setting.Config{}
	// err := c.Get()
	// if err != nil {
	// 	log.Println("Failed to initialize logger: &v", err)
	// }
	f, _ := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	Log = log.New(f, "", log.Lshortfile|log.Ltime)
}
