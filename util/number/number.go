package number

import (
	"math"
	"regexp"
)

func RoundFloat[T float32 | float64](val T, precision uint) T {
	ratio := math.Pow(10, float64(precision))
	return T(math.Round(float64(val)*ratio) / ratio)
}

var notNumberRE = regexp.MustCompile(`\D`)

// CleanNumber remove all not number characters from string
func CleanNumber(value string) string {
	return notNumberRE.ReplaceAllString(value, "")
}
