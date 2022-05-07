package crypto

import (
	"crypto/md5"
	"encoding/hex"
)

// GetMd5 - to encrypt some string
func GetMd5(input string) string {
	//TODO change how to generate password more strong
	hash := md5.New()
	defer hash.Reset()
	hash.Write([]byte(input))

	return hex.EncodeToString(hash.Sum(nil))
}
