package checksumutils

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
)

func GetSha1SumFromBytes(bytesToHash []byte) (checksum string) {
	// Source: https://gobyexample.com/sha256-hashes
	h := sha1.New()
	h.Write(bytesToHash)
	checksumBytes := h.Sum(nil)
	checksum = fmt.Sprintf("%x", checksumBytes)
	return checksum
}

func GetSha1SumFromString(stringToHash string) (checksum string) {
	return GetSha1SumFromBytes([]byte(stringToHash))
}

func GetSha256SumFromBytes(bytesToHash []byte) (checksum string) {
	h := sha256.New()
	h.Write(bytesToHash)
	checksumBytes := h.Sum(nil)
	checksum = fmt.Sprintf("%x", checksumBytes)
	return checksum
}

func GetSha256SumFromString(stringToHash string) (checksum string) {
	return GetSha256SumFromBytes([]byte(stringToHash))
}

func GetSha512SumFromBytes(bytesToHash []byte) (checksum string) {
	h := sha512.New()
	h.Write(bytesToHash)
	checksumBytes := h.Sum(nil)
	checksum = fmt.Sprintf("%x", checksumBytes)
	return checksum
}

func GetSha512SumFromString(stringToHash string) (checksum string) {
	return GetSha512SumFromBytes([]byte(stringToHash))
}
