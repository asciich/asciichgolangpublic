package asciichgolangpublic

import (
	"os"
)

// A File represents any kind of file regardless if a local file or a remote file.
type File interface {
	Create(verbose bool) (err error)
	Delete(verbose bool) (err error)
	Exists() (exists bool, err error)
	GetBaseName() (baseName string, err error)
	GetLocalPath() (localPath string, err error)
	GetUriAsString() (uri string, err error)
	MustCreate(verbose bool)
	MustDelete(verbose bool)
	MustExists() (exists bool)
	MustGetBaseName() (baseName string)
	MustGetLocalPath() (localPath string)
	MustGetUriAsString() (uri string)
	MustReadAsBytes() (content []byte)
	MustWriteBytes(toWrite []byte, verbose bool)
	ReadAsBytes() (content []byte, err error)
	WriteBytes(toWrite []byte, verbose bool) (err error)

	// All methods below this line can be implemented by embedding the `FileBase` struct:
	GetSha256Sum() (sha256sum string, err error)
	IsMatchingSha256Sum(sha256sum string) (isMatching bool, err error)
	MustGetSha256Sum() (sha256sum string)
	MustIsMatchingSha256Sum(sha256sum string) (isMatching bool)
	MustReadAsString() (content string)
	MustWriteString(content string, verbose bool)
	ReadAsString() (content string, err error)
	WriteString(content string, verbose bool) (err error)
}

func GetFileByOsFile(osFile *os.File) (file File, err error) {
	if osFile == nil {
		return nil, TracedError("osFile is nil")
	}

	file, err = NewLocalFileByPath(osFile.Name())
	if err != nil {
		return nil, err
	}

	return file, nil
}

func MustGetFileByOsFile(osFile *os.File) (file File) {
	file, err := GetFileByOsFile(osFile)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return file
}
