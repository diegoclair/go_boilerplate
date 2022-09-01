package number

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFloatingPoint(t *testing.T) {

	x := 0.07
	x1 := 0.07
	x2 := 0.07

	totalx := x + x1 + x2

	require.Less(t, 0.21, totalx)                 //different value because of floating point
	require.Equal(t, 0.21, RoundFloat(totalx, 2)) //fixed floating point value

	y := 0.1
	y1 := 0.2
	totaly := y + y1

	require.Less(t, 0.3, totaly)                 //different value because of floating point
	require.Equal(t, 0.3, RoundFloat(totaly, 2)) //fixed floating point value
}
