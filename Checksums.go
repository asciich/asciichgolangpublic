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

func (c *ChecksumsService) GetSha256SumFromBytes(bytesToHash []byte) (checksum string) {
	// Source: https://gobyexample.com/sha256-hashes
	h := sha256.New()
	h.Write(bytesToHash)
	checksumBytes := h.Sum(nil)
	checksum = fmt.Sprintf("%x", checksumBytes)
	return checksum
}

func (c *ChecksumsService) GetSha256SumFromString(stringToHash string) (checksum string) {
	return c.GetSha256SumFromBytes([]byte(stringToHash))
}
