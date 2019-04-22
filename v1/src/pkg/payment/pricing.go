package payment

import (
	"encoding/json"
	"os"
)

// ride cost constants are stored in a file
func filename() string { return "pricing.json" }

// Price initializes at start time
var (
	Price      Pricing
	OtherPrice Pricing
)

// Pricing has the constants needed to calculate prices
type Pricing struct {
	BaseFare float64 `json:"base_fare"`
	Time     float64 `json:"time"`     // time in minutes
	Distance float64 `json:"distance"` // distance in kilometer
}

func init() {
	Price.Load("nduthi_pricing.json")
	OtherPrice.Load("car_pricing.json")
}

// Load gets the values for pricing
func (p *Pricing) Load(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&p)
	if err != nil {
		return err
	}
	return nil
}

// UpdateJSON updates pricing constants in the file
func (p *Pricing) UpdateJSON() error {
	f, err := os.OpenFile(filename(), os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(&p)
	if err != nil {
		return err
	}
	return nil
}

// Vehicle type distinct pricing
func Vehicle(vType string) *Pricing {
	switch vType {
	case "nduthi":
		return &Price
	case "car":
		return &OtherPrice
	}
	return nil
}
