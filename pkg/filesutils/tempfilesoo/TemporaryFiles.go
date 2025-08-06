package tempfilesoo

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func CreateEmptyTemporaryFile(ctx context.Context) (temporaryfile filesinterfaces.File, err error) {
	temporaryfile, err = CreateNamedTemporaryFile(ctx, "emptyFile")
	if err != nil {
		return nil, err
	}

	return temporaryfile, nil
}

func CreateEmptyTemporaryFileAndGetPath(ctx context.Context) (temporaryFilePath string, err error) {
	temporaryFile, err := CreateEmptyTemporaryFile(ctx)
	if err != nil {
		return "", err
	}

	temporaryFilePath, err = temporaryFile.GetLocalPath()
	if err != nil {
		return "", err
	}

	return temporaryFilePath, nil
}

func CreateFromBytes(ctx context.Context, content []byte) (temporaryFile filesinterfaces.File, err error) {
	if content == nil {
		return nil, tracederrors.TracedErrorNil("content")
	}

	temporaryFile, err = CreateEmptyTemporaryFile(ctx)
	if err != nil {
		return nil, err
	}

	err = temporaryFile.WriteBytes(ctx, content, &filesoptions.WriteOptions{})
	if err != nil {
		return nil, err
	}

	return temporaryFile, nil
}

func CreateFromString(ctx context.Context, content string) (temporaryFile filesinterfaces.File, err error) {
	temporaryFile, err = CreateEmptyTemporaryFile(ctx)
	if err != nil {
		return nil, err
	}

	err = temporaryFile.WriteString(ctx, content, &filesoptions.WriteOptions{})
	if err != nil {
		return nil, err
	}

	return temporaryFile, nil
}

func CreateFromStringAndGetPath(ctx context.Context, content string) (temporaryFilePath string, err error) {
	temporaryFile, err := CreateFromString(ctx, content)
	if err != nil {
		return "", err
	}

	temporaryFilePath, err = temporaryFile.GetPath()
	if err != nil {
		return "", err
	}

	return temporaryFilePath, nil
}

func CreateNamedTemporaryFile(ctx context.Context, fileName string) (temporaryfile filesinterfaces.File, err error) {
	if fileName == "" {
		return nil, tracederrors.TracedErrorEmptyString("fileName")
	}

	tmpPath, err := tempfiles.CreateNamedTemporaryFile(ctx, fileName)
	if err != nil {
		return nil, err
	}

	temporaryfile, err = files.GetLocalFileByPath(tmpPath)
	if err != nil {
		return nil, err
	}

	createdFilePath, err := temporaryfile.GetPath()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Created temporary file '%s'", createdFilePath)

	return temporaryfile, nil
}

func CreateTemporaryFileFromBytes(ctx context.Context, content []byte) (temporaryFile filesinterfaces.File, err error) {
	if content == nil {
		return nil, tracederrors.TracedErrorNil("content")
	}

	temporaryFile, err = CreateTemporaryFileFromString(ctx, string(content))
	if err != nil {
		return nil, err
	}

	return temporaryFile, nil
}

func CreateTemporaryFileFromFile(ctx context.Context, fileToCopyAsTemporaryFile filesinterfaces.File) (temporaryFile filesinterfaces.File, err error) {
	if fileToCopyAsTemporaryFile == nil {
		return nil, tracederrors.TracedErrorNil("fileToCopyAsTemporaryFile")
	}

	hostDescription, err := fileToCopyAsTemporaryFile.GetHostDescription()
	if err != nil {
		return nil, err
	}

	isLocalFile, err := fileToCopyAsTemporaryFile.IsLocalFile(contextutils.GetVerboseFromContext(ctx))
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

	temporaryFile, err = CreateNamedTemporaryFile(ctx, filepath.Base(fileToCopyAsTemporaryFilePath)+"_tmp")
	if err != nil {
		return nil, err
	}

	err = fileToCopyAsTemporaryFile.CopyToFile(temporaryFile, contextutils.GetVerboseFromContext(ctx))
	if err != nil {
		return nil, err
	}

	temporaryfilePath, err := temporaryFile.GetPath()
	if err != nil {
		return nil, err
	}

	srcPath, err := fileToCopyAsTemporaryFile.GetPath()
	if err != nil {
		return nil, err
	}

	logging.LogChangedByCtxf(
		ctx,
		"Created temporary file '%s' filed with content from '%s' on '%s'",
		temporaryfilePath,
		srcPath,
		hostDescription,
	)

	return temporaryFile, nil
}

func CreateTemporaryFileFromPath(ctx context.Context, filePathToCopyAsTemporaryFile ...string) (temporaryFile filesinterfaces.File, err error) {
	if len(filePathToCopyAsTemporaryFile) <= 0 {
		return nil, tracederrors.TracedError("filePathToCopyAsTemporaryFile")
	}

	pathToUse := strings.Join(filePathToCopyAsTemporaryFile, "/")

	fileToCopy, err := files.GetLocalFileByPath(pathToUse)
	if err != nil {
		return nil, err
	}

	temporaryFile, err = CreateTemporaryFileFromFile(ctx, fileToCopy)
	if err != nil {
		return nil, err
	}

	return temporaryFile, nil
}

func CreateTemporaryFileFromString(ctx context.Context, content string) (temporaryFile filesinterfaces.File, err error) {
	temporaryFile, err = CreateNamedTemporaryFile(ctx, "tempFile")
	if err != nil {
		return nil, err
	}

	err = temporaryFile.WriteString(ctx, content, &filesoptions.WriteOptions{})
	if err != nil {
		return nil, err
	}

	return temporaryFile, nil
}
