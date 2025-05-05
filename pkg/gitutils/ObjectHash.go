package gitutils

import (
	"bytes"
	"fmt"

	"github.com/asciich/asciichgolangpublic/checksums"
)

func GetBlobOjectHashFromString(content string) (hash string) {
	return GetBlobObjectHashFromBytes([]byte(content))
}

func GetBlobObjectHashFromBytes(content []byte) (hash string) {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("blob %d\x00", len(content)))
	buf.Write(content)

	return checksums.GetSha1SumFromBytes(buf.Bytes())
}
