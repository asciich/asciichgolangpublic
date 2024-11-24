package asciichgolangpublic

type FilesService struct {
}

func Files() (f *FilesService) {
	return NewFilesService()
}

func MustReadFileAsString(path string) (content string) {
	content, err := ReadFileAsString(path)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return content
}

func MustWriteStringToFile(path string, content string, verbose bool) {
	err := WriteStringToFile(path, content, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func NewFilesService() (f *FilesService) {
	return new(FilesService)
}

func ReadFileAsString(path string) (content string, err error) {
	return Files().ReadAsString(path)
}

func WriteStringToFile(path string, content string, verbose bool) (err error) {
	return Files().WriteStringToFile(path, content, verbose)
}

func (f *FilesService) MustReadAsString(path string) (content string) {
	content, err := f.ReadAsString(path)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return content
}

func (f *FilesService) MustWriteStringToFile(path string, content string, verbose bool) {
	err := f.WriteStringToFile(path, content, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (f *FilesService) ReadAsString(path string) (content string, err error) {
	if path == "" {
		return "", TracedErrorEmptyString(path)
	}

	localFile, err := GetLocalFileByPath(path)
	if err != nil {
		return "", err
	}

	content, err = localFile.ReadAsString()
	if err != nil {
		return "", err
	}

	return content, nil
}

func (f *FilesService) WriteStringToFile(path string, content string, verbose bool) (err error) {
	if path == "" {
		return TracedErrorEmptyString(path)
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
