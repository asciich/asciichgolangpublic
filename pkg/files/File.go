package files

import (
	"os"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func GetFileByOsFile(osFile *os.File) (file filesinterfaces.File, err error) {
	if osFile == nil {
		return nil, tracederrors.TracedError("osFile is nil")
	}

	file, err = NewLocalFileByPath(osFile.Name())
	if err != nil {
		return nil, err
	}

	return file, nil
}

func MustGetFileByOsFile(osFile *os.File) (file filesinterfaces.File) {
	file, err := GetFileByOsFile(osFile)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return file
}
