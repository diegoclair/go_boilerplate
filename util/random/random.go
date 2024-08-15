package random

import (
	"math/rand"
	"strings"
	"time"

	"github.com/mvrilo/go-cpf"
)

const (
	alphabet       = "abcdefghijklmnopqrstuvwxyz"
	alphabetLength = len(alphabet)
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomCPF() string {
	return cpf.Generate()
}

func RandomName() string {
	return RandomString(6)
}

func RandomPassword() string {
	return RandomString(8)
}

// RandomString generates a random string of length n.
// It uses the characters from the alphabet defined in the package.
// The generated string is returned as a result.
func RandomString(n int) string {
	var sb strings.Builder

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(alphabetLength)]
		sb.WriteByte(c)
	}
	return sb.String()
}
