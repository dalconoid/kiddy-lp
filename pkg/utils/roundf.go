package utils

import "math"

// RoundF64ToPrecision rounds float f to precision p
func RoundF64ToPrecision(f float64, p int) float64 {
	mult := math.Pow10(p)
	return math.Round(f * mult) / mult
}