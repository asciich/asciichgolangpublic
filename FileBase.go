package asciichgolangpublic

// This is the base for `File` providing most convenience functions for file operations.
type FileBase struct {
	parentFileForBaseClass File
}

func NewFileBase() (f *FileBase) {
	return new(FileBase)
}

func (f *FileBase) GetParentFileForBaseClass() (parentFileForBaseClass File, err error) {

	return f.parentFileForBaseClass, nil
}

func (f *FileBase) MustGetParentFileForBaseClass() (parentFileForBaseClass File) {
	parentFileForBaseClass, err := f.GetParentFileForBaseClass()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return parentFileForBaseClass
}

func (f *FileBase) MustReadAsString() (content string) {
	content, err := f.ReadAsString()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return content
}

func (f *FileBase) MustSetParentFileForBaseClass(parentFileForBaseClass File) {
	err := f.SetParentFileForBaseClass(parentFileForBaseClass)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (f *FileBase) MustWriteString(toWrite string, verbose bool) {
	err := f.WriteString(toWrite, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (f *FileBase) ReadAsString() (content string, err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return "", err
	}

	contentBytes, err := parent.ReadAsBytes()
	if err != nil {
		return "", err
	}

	return string(contentBytes), nil
}

func (f *FileBase) SetParentFileForBaseClass(parentFileForBaseClass File) (err error) {
	f.parentFileForBaseClass = parentFileForBaseClass

	return nil
}

func (f *FileBase) WriteString(toWrite string, verbose bool) (err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return err
	}

	return parent.WriteBytes([]byte(toWrite), verbose)
}
