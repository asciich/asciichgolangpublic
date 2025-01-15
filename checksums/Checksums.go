package checksums

import (
	"crypto/sha256"
	"fmt"
)

func GetSha256SumFromBytes(bytesToHash []byte) (checksum string) {
	// Source: https://gobyexample.com/sha256-hashes
	h := sha256.New()
	h.Write(bytesToHash)
	checksumBytes := h.Sum(nil)
	checksum = fmt.Sprintf("%x", checksumBytes)
	return checksum
}

func GetSha256SumFromString(stringToHash string) (checksum string) {
	return GetSha256SumFromBytes([]byte(stringToHash))
}
