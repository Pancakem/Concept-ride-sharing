package payment

import (
	"encoding/json"
	"os"
)

// ride cost constants are stored in a file
const filename = "pricing.json"

// Pricing has the constants needed to calculate prices
type Pricing struct {
	BaseFare float64 `json:"base_fare"`
	Time     float64 `json:"time"`
	Distance float64 `json:"distance"`
}

// Load gets the values for pricing
func (p *Pricing) Load() error {
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
	f, err := os.OpenFile(filename, os.O_WRONLY, 0666)
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
