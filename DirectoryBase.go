package asciichgolangpublic

type DirectoryBase struct {
	parentDirectoryForBaseClass Directory
}

func NewDirectoryBase() (d *DirectoryBase) {
	return new(DirectoryBase)
}

func (d *DirectoryBase) GetFilePathInDirectory(path ...string) (filePath string, err error) {
	if len(path) <= 0 {
		return "", TracedError("path has no elements")
	}

	parent, err := d.GetParentDirectoryForBaseClass()
	if err != nil {
		return "", nil
	}

	localFile, err := parent.GetFileInDirectory(path...)
	if err != nil {
		return "", err
	}

	filePath, err = localFile.GetLocalPath()
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func (d *DirectoryBase) GetParentDirectoryForBaseClass() (parentDirectoryForBaseClass Directory, err error) {
	if d.parentDirectoryForBaseClass == nil {
		return nil, TracedError("parentDirectoryForBaseClass not set")
	}
	return d.parentDirectoryForBaseClass, nil
}

func (d *DirectoryBase) MustGetFilePathInDirectory(path ...string) (filePath string) {
	filePath, err := d.GetFilePathInDirectory(path...)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return filePath
}

func (d *DirectoryBase) MustGetParentDirectoryForBaseClass() (parentDirectoryForBaseClass Directory) {
	parentDirectoryForBaseClass, err := d.GetParentDirectoryForBaseClass()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return parentDirectoryForBaseClass
}

func (d *DirectoryBase) MustSetParentDirectoryForBaseClass(parentDirectoryForBaseClass Directory) {
	err := d.SetParentDirectoryForBaseClass(parentDirectoryForBaseClass)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (d *DirectoryBase) SetParentDirectoryForBaseClass(parentDirectoryForBaseClass Directory) (err error) {
	if parentDirectoryForBaseClass == nil {
		return TracedErrorNil("parentDirectoryForBaseClass")
	}

	d.parentDirectoryForBaseClass = parentDirectoryForBaseClass

	return nil
}
