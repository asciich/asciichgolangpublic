package gitutils

import (
	"bytes"
	"fmt"

	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/checksumutils"
)

func GetBlobOjectHashFromString(content string) (hash string) {
	return GetBlobObjectHashFromBytes([]byte(content))
}

func GetBlobObjectHashFromBytes(content []byte) (hash string) {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("blob %d\x00", len(content)))
	buf.Write(content)

	return checksumutils.GetSha1SumFromBytes(buf.Bytes())
}
