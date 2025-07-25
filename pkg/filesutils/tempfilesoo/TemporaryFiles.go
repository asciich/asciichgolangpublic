package tempfilesoo

import (
	"path/filepath"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func CreateEmptyTemporaryFile(verbose bool) (temporaryfile filesinterfaces.File, err error) {
	temporaryfile, err = CreateNamedTemporaryFile("emptyFile", verbose)
	if err != nil {
		return nil, err
	}

	return temporaryfile, nil
}

func CreateEmptyTemporaryFileAndGetPath(verbose bool) (temporaryFilePath string, err error) {
	temporaryFile, err := CreateEmptyTemporaryFile(verbose)
	if err != nil {
		return "", err
	}

	temporaryFilePath, err = temporaryFile.GetLocalPath()
	if err != nil {
		return "", err
	}

	return temporaryFilePath, nil
}

func CreateFromBytes(content []byte, verbose bool) (temporaryFile filesinterfaces.File, err error) {
	if content == nil {
		return nil, tracederrors.TracedErrorNil("content")
	}

	temporaryFile, err = CreateEmptyTemporaryFile(verbose)
	if err != nil {
		return nil, err
	}

	err = temporaryFile.WriteBytes(content, verbose)
	if err != nil {
		return nil, err
	}

	return temporaryFile, nil
}

func CreateFromString(content string, verbose bool) (temporaryFile filesinterfaces.File, err error) {
	temporaryFile, err = CreateEmptyTemporaryFile(verbose)
	if err != nil {
		return nil, err
	}

	err = temporaryFile.WriteString(content, verbose)
	if err != nil {
		return nil, err
	}

	return temporaryFile, nil
}

func CreateFromStringAndGetPath(content string, verbose bool) (temporaryFilePath string, err error) {
	temporaryFile, err := CreateFromString(content, verbose)
	if err != nil {
		return "", err
	}

	temporaryFilePath, err = temporaryFile.GetPath()
	if err != nil {
		return "", err
	}

	return temporaryFilePath, nil
}

func CreateNamedTemporaryFile(fileName string, verbose bool) (temporaryfile filesinterfaces.File, err error) {
	if fileName == "" {
		return nil, tracederrors.TracedErrorEmptyString("fileName")
	}

	tmpPath, err := tempfiles.CreateNamedTemporaryFile(contextutils.GetVerbosityContextByBool(verbose), fileName)
	if err != nil {
		return nil, err
	}

	temporaryfile, err = files.GetLocalFileByPath(tmpPath)
	if err != nil {
		return nil, err
	}

	if verbose {
		createdFilePath, err := temporaryfile.GetPath()
		if err != nil {
			return nil, err
		}

		logging.LogInfof("Created temporary file '%s'", createdFilePath)
	}

	return temporaryfile, nil
}

func CreateTemporaryFileFromBytes(content []byte, verbose bool) (temporaryFile filesinterfaces.File, err error) {
	if content == nil {
		return nil, tracederrors.TracedErrorNil("content")
	}

	temporaryFile, err = CreateTemporaryFileFromString(string(content), verbose)
	if err != nil {
		return nil, err
	}

	return temporaryFile, nil
}

func CreateTemporaryFileFromFile(fileToCopyAsTemporaryFile filesinterfaces.File, verbose bool) (temporaryFile filesinterfaces.File, err error) {
	if fileToCopyAsTemporaryFile == nil {
		return nil, tracederrors.TracedErrorNil("fileToCopyAsTemporaryFile")
	}

	hostDescription, err := fileToCopyAsTemporaryFile.GetHostDescription()
	if err != nil {
		return nil, err
	}

	isLocalFile, err := fileToCopyAsTemporaryFile.IsLocalFile(verbose)
	if err != nil {
		return nil, err
	}

	if !isLocalFile {
		return nil, tracederrors.TracedErrorf(
			"Only implemented for files on 'localhost' but got '%s'",
			hostDescription,
		)
	}

	fileToCopyAsTemporaryFilePath, err := fileToCopyAsTemporaryFile.GetPath()
	if err != nil {
		return nil, err
	}

	temporaryFile, err = CreateNamedTemporaryFile(filepath.Base(fileToCopyAsTemporaryFilePath)+"_tmp", verbose)
	if err != nil {
		return nil, err
	}

	err = fileToCopyAsTemporaryFile.CopyToFile(temporaryFile, verbose)
	if err != nil {
		return nil, err
	}

	if verbose {
		temporaryfilePath, err := temporaryFile.GetPath()
		if err != nil {
			return nil, err
		}

		srcPath, err := fileToCopyAsTemporaryFile.GetPath()
		if err != nil {
			return nil, err
		}

		logging.LogChangedf(
			"Created temporary file '%s' filed with content from '%s' on '%s'",
			temporaryfilePath,
			srcPath,
			hostDescription,
		)
	}

	return temporaryFile, nil
}

func CreateTemporaryFileFromPath(verbose bool, filePathToCopyAsTemporaryFile ...string) (temporaryFile filesinterfaces.File, err error) {
	if len(filePathToCopyAsTemporaryFile) <= 0 {
		return nil, tracederrors.TracedError("filePathToCopyAsTemporaryFile")
	}

	pathToUse := strings.Join(filePathToCopyAsTemporaryFile, "/")

	fileToCopy, err := files.GetLocalFileByPath(pathToUse)
	if err != nil {
		return nil, err
	}

	temporaryFile, err = CreateTemporaryFileFromFile(fileToCopy, verbose)
	if err != nil {
		return nil, err
	}

	return temporaryFile, nil
}

func CreateTemporaryFileFromString(content string, verbose bool) (temporaryFile filesinterfaces.File, err error) {
	temporaryFile, err = CreateNamedTemporaryFile("tempFile", verbose)
	if err != nil {
		return nil, err
	}

	err = temporaryFile.WriteString(content, verbose)
	if err != nil {
		return nil, err
	}

	return temporaryFile, nil
}
