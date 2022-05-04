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
	rand.Seed(time.Now().UnixNano())
}

func RandomCPF() string {
	return cpf.Generate()
}

func RandomName() string {
	return RandomString(6)
}

func RandomSecret() string {
	return RandomString(8)
}

func RandomString(n int) string {
	var sb strings.Builder

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(alphabetLength)]
		sb.WriteByte(c)
	}
	return sb.String()
}
