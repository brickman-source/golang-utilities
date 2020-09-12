package crypto

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"golang.org/x/crypto/sha3"
)

func GetSha1Hash(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func GetSha256Hash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func GetSha384Hash(s string) string {
	h := sha3.New384()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
