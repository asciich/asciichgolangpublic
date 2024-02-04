package asciichgolangpublic

import (
	"crypto/sha256"
	"fmt"
)

type ChecksumsService struct {
}

func Checksums() (checksums *ChecksumsService) {
	return new(ChecksumsService)
}

func NewChecksumsService() (c *ChecksumsService) {
	return new(ChecksumsService)
}

func (c *ChecksumsService) GetSha256SumFromString(stringToHash string) (checksum string) {
	// Source: https://gobyexample.com/sha256-hashes
	h := sha256.New()
	h.Write([]byte(stringToHash))
	checksumBytes := h.Sum(nil)
	checksum = fmt.Sprintf("%x", checksumBytes)
	return checksum
}
