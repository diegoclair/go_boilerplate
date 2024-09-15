package number

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFloatingPointFloat64(t *testing.T) {
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

func TestFloatingPointFloat32(t *testing.T) {
	x := float32(0.3)
	x1 := float32(0.1)
	totalx := x - x1

	require.Less(t, float32(0.2), totalx)                 //different value because of floating point
	require.Equal(t, float32(0.20000002), totalx)         //floating point problem
	require.Equal(t, float32(0.2), RoundFloat(totalx, 2)) //fixed floating point value

	y := float32(0.1)
	y1 := float32(0.2)
	totaly := y + y1

	require.Equal(t, float32(0.3), totaly)                // not floating point problem in this case
	require.Equal(t, float32(0.3), RoundFloat(totaly, 2)) //fixed floating point value

	z := float32(0.07)
	z1 := float32(0.07)
	z2 := float32(0.07)

	totalz := z + z1 + z2
	require.Less(t, float32(0.21), totalz)                 //different value because of floating point
	require.Equal(t, float32(0.21000001), totalz)          //floating point problem
	require.Equal(t, float32(0.21), RoundFloat(totalz, 2)) //fixed floating point value

	tax := totalz * float32(0.05)
	require.Less(t, float32(0.0105), tax)                         //different value because of floating point
	require.Equal(t, float32(0.010500001), tax)                   //floating point problem
	require.Equal(t, float32(0.0105), RoundFloat(totalz, 4)*0.05) //fixed floating point value

	final := totalz + tax
	require.Equal(t, float32(0.22050001), final)            //floating point problem
	require.Equal(t, float32(0.2205), RoundFloat(final, 4)) //fixed floating point value
}

func TestFloatingPointEdgeCases(t *testing.T) {
	t.Run("float64 edge cases", func(t *testing.T) {
		// Subtraction leading to small imprecision
		x := 0.3
		x1 := 0.1
		totalx := x - x1
		require.NotEqual(t, 0.2, totalx)
		require.Equal(t, 0.2, RoundFloat(totalx, 2))

		// Addition leading to small imprecision
		y := 0.1
		y1 := 0.2
		totaly := y + y1
		require.NotEqual(t, 0.3, totaly)
		require.Equal(t, 0.3, RoundFloat(totaly, 2))

		// Multiplication leading to imprecision
		z := 0.1
		z1 := 3.0
		totalz := z * z1
		require.NotEqual(t, 0.3, totalz)
		require.Equal(t, 0.3, RoundFloat(totalz, 2))

		// Division leading to repeating decimal
		w := 1.0
		w1 := 3.0
		totalw := w / w1
		require.NotEqual(t, 0.33, totalw)
		require.Equal(t, 0.33, RoundFloat(totalw, 2))
	})

	t.Run("float32 edge cases", func(t *testing.T) {
		// Subtraction leading to small imprecision
		x := float32(0.3)
		x1 := float32(0.1)
		totalx := x - x1
		require.NotEqual(t, float32(0.2), totalx)
		require.Equal(t, float32(0.2), RoundFloat(totalx, 2))

		// Addition leading to small imprecision
		y := float32(0.1)
		y1 := float32(0.2)
		totaly := y + y1
		// Note: This might actually be equal due to float32's lower precision
		require.Equal(t, float32(0.3), totaly)

		z := float32(0.1234)
		z1 := float32(3)
		totalz := z * z1
		require.NotEqual(t, float32(0.37), totalz)             // totalz should be about 0.3702
		require.Equal(t, float32(0.37), RoundFloat(totalz, 2)) // This checks your rounding function

		// Division leading to repeating decimal
		w := float32(1.0)
		w1 := float32(3.0)
		totalw := w / w1
		require.NotEqual(t, float32(0.33), totalw)
		require.Equal(t, float32(0.33), RoundFloat(totalw, 2))

		z = float32(0.1)
		z1 = float32(3)
		totalz = z * z1
		require.Equal(t, float32(0.3), totalz)                 // This checks Go's float32 behavior
		require.Equal(t, float32(0.30), RoundFloat(totalz, 2)) // This checks your rounding function
	})

	t.Run("complex calculations", func(t *testing.T) {
		// A series of operations that accumulate error
		a := 0.1
		b := 0.2
		c := 0.3
		result := (a + b + c) * 10 / 3
		require.NotEqual(t, 2.0, result)
		require.Equal(t, 2.0, RoundFloat(result, 2))

		// Same calculation with float32
		a32 := float32(0.1)
		b32 := float32(0.2)
		c32 := float32(0.3)
		result32 := (a32 + b32 + c32) * 10 / 3
		require.Equal(t, float32(2.0), result32)
		require.Equal(t, float32(2.0), RoundFloat(result32, 2))
	})
}

func TestCleanNumber(t *testing.T) {
	require.Equal(t, "1234567890", CleanNumber("1234567890"))
	require.Equal(t, "1234567890", CleanNumber("1234567890a"))
	require.Equal(t, "1234567890", CleanNumber("a1234567890"))
	require.Equal(t, "1234567890", CleanNumber("a1234567890a"))
	require.Equal(t, "1234567890", CleanNumber("a1b2c3d4e5f6g7h8i9j0"))
	require.Equal(t, "1234567890", CleanNumber("a1.b2-c3*d4/e5+f6=g7h8i9j0"))
}
