package number

import (
	"math"
	"regexp"
)

func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

var notNumberRE = regexp.MustCompile(`\D`)

// CleanNumber remove all not number characters from string
func CleanNumber(value string) string {
	return notNumberRE.ReplaceAllString(value, "")
}
