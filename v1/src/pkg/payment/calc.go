package payment

// Calculate returns the estimated fare
// it is also used in the actual fare calculation
func (p *Pricing) Calculate(time, distance float64) float64 {
	return p.BaseFare + (p.Time * time) + (p.Distance * distance)
}
