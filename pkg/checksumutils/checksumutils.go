package checksumutils

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
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

func GetMD5SumFromBytes(bytesToHash []byte) (checksum string) {
	h := md5.New()
	h.Write(bytesToHash)
	checksumBytes := h.Sum(nil)
	checksum = fmt.Sprintf("%x", checksumBytes)
	return checksum
}

func GetMD5SumFromString(stringToHash string) (checksum string) {
	return GetMD5SumFromBytes([]byte(stringToHash))
}

func GetMD5SumFromFileByPath(ctx context.Context, path string) (checksum string, err error) {
	if path == "" {
		return "", tracederrors.TracedErrorEmptyString("path")
	}

	logging.LogInfoByCtxf(ctx, "Get MD5 checksum of '%s' started.", path)

	file, err := os.Open(path)
	if err != nil {
		return "", tracederrors.TracedErrorf("Failed to open '%s': %w", path, err)
	}
	defer file.Close()

	hash := md5.New()

	// Use a channel to handle the hashing in case the context is cancelled
	// for very large files, though io.Copy is generally fast.
	errChan := make(chan error, 1)

	go func() {
		// io.Copy streams the file content into the hash object
		_, copyErr := io.Copy(hash, file)
		errChan <- copyErr
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case err := <-errChan:
		if err != nil {
			return "", fmt.Errorf("failed to hash file: %w", err)
		}
	}

	checksum = hex.EncodeToString(hash.Sum(nil))

	logging.LogInfoByCtxf(ctx, "Get MD5 checksum of '%s' finished. Checksum is '%s'.", path, checksum)

	return checksum, nil
}
