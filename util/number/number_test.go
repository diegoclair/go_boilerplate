package number

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFloatingPoint(t *testing.T) {

	x := 0.3
	x1 := 0.1
	totalx := x - x1

	require.Greater(t, 0.2, totalx)               //different value because of floating point
	require.Equal(t, 0.19999999999999998, totalx) //floating point problem
	require.Equal(t, 0.2, RoundFloat(totalx, 2))  //fixed floating point value

	y := 0.1
	y1 := 0.2
	totaly := y + y1

	require.Less(t, 0.3, totaly)                  //different value because of floating point
	require.Equal(t, 0.30000000000000004, totaly) //floating point problem
	require.Equal(t, 0.3, RoundFloat(totaly, 2))  //fixed floating point value

	z := 0.07
	z1 := 0.07
	z2 := 0.07

	totalz := z + z1 + z2
	require.Less(t, 0.21, totalz)                 //different value because of floating point
	require.Equal(t, 0.21000000000000002, totalz) //floating point problem
	require.Equal(t, 0.21, RoundFloat(totalz, 2)) //fixed floating point value

	tax := totalz * 0.05
	require.Less(t, 0.0105, tax)                         //different value because of floating point
	require.Equal(t, 0.010500000000000002, tax)          //floating point problem
	require.Equal(t, 0.0105, RoundFloat(totalz, 4)*0.05) //fixed floating point value

	final := totalz + tax
	require.Less(t, 0.2205, final)                 //different value because of floating point
	require.Equal(t, 0.22050000000000003, final)   //floating point problem
	require.Equal(t, 0.2205, RoundFloat(final, 4)) //fixed floating point value
}

func TestCleanNumber(t *testing.T) {
	require.Equal(t, "1234567890", CleanNumber("1234567890"))
	require.Equal(t, "1234567890", CleanNumber("1234567890a"))
	require.Equal(t, "1234567890", CleanNumber("a1234567890"))
	require.Equal(t, "1234567890", CleanNumber("a1234567890a"))
	require.Equal(t, "1234567890", CleanNumber("a1b2c3d4e5f6g7h8i9j0"))
	require.Equal(t, "1234567890", CleanNumber("a1.b2-c3*d4/e5+f6=g7h8i9j0"))
}
