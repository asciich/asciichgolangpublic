package asciichgolangpublic

type FilesService struct {
}

func Files() (f *FilesService) {
	return NewFilesService()
}

func NewFilesService() (f *FilesService) {
	return new(FilesService)
}

func (f *FilesService) MustWriteStringToFile(path string, content string, verbose bool) {
	err := f.WriteStringToFile(path, content, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (f *FilesService) WriteStringToFile(path string, content string, verbose bool) (err error) {
	if path == "" {
		return TracedErrorNil(path)
	}

	localFile, err := GetLocalFileByPath(path)
	if err != nil {
		return err
	}

	err = localFile.WriteString(content, verbose)
	if err != nil {
		return err
	}

	return nil
}
