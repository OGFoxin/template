package utils

import (
	"fmt"
	"math"
)

func RoundTo(num float64, n int) string {
	pow := math.Pow(10, float64(n))
	return fmt.Sprintf("%.2f", math.Round(num*pow)/pow)
}
