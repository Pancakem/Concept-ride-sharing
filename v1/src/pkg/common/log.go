package common

import (
	"log"
	"os"
)

var (
	// Log write logs to file
	Log *log.Logger
)

func init() {
	f, _ := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	Log = log.New(f, "", log.Lshortfile|log.Ltime)
}
