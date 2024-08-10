package asciichgolangpublic

type DirectoryBase struct {
	parentDirectoryForBaseClass Directory
}

func NewDirectoryBase() (d *DirectoryBase) {
	return new(DirectoryBase)
}

func (d *DirectoryBase) CreateFileInDirectoryFromString(content string, verbose bool, pathToCreate ...string) (createdFile File, err error) {
	if len(pathToCreate) <= 0 {
		return nil, TracedErrorf("Invalid pathToCreate='%v'", pathToCreate)
	}

	parent, err := d.GetParentDirectoryForBaseClass()
	if err != nil {
		return nil, err
	}

	createdFile, err = parent.GetFileInDirectory(pathToCreate...)
	if err != nil {
		return nil, err
	}

	parentDir, err := createdFile.GetParentDirectory()
	if err != nil {
		return nil, err
	}

	err = parentDir.Create(verbose)
	if err != nil {
		return nil, err
	}

	err = createdFile.WriteString(content, verbose)
	if err != nil {
		return nil, err
	}

	return createdFile, nil
}

func (d *DirectoryBase) FileInDirectoryExists(path ...string) (fileExists bool, err error) {
	if len(path) <= 0 {
		return false, TracedError("path has no elements")
	}

	parent, err := d.GetParentDirectoryForBaseClass()
	if err != nil {
		return false, err
	}

	fileToCheck, err := parent.GetFileInDirectory(path...)
	if err != nil {
		return false, err
	}

	fileExists, err = fileToCheck.Exists()
	if err != nil {
		return false, err
	}

	return fileExists, nil
}

func (d *DirectoryBase) GetFilePathInDirectory(path ...string) (filePath string, err error) {
	if len(path) <= 0 {
		return "", TracedError("path has no elements")
	}

	parent, err := d.GetParentDirectoryForBaseClass()
	if err != nil {
		return "", err
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

func (d *DirectoryBase) MustCreateFileInDirectoryFromString(content string, verbose bool, pathToCreate ...string) (createdFile File) {
	createdFile, err := d.CreateFileInDirectoryFromString(content, verbose, pathToCreate...)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return createdFile
}

func (d *DirectoryBase) MustFileInDirectoryExists(path ...string) (fileExists bool) {
	fileExists, err := d.FileInDirectoryExists(path...)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return fileExists
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

func (d *DirectoryBase) MustReadFileInDirectoryAsLines(path ...string) (content []string) {
	content, err := d.ReadFileInDirectoryAsLines(path...)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return content
}

func (d *DirectoryBase) MustReadFileInDirectoryAsString(path ...string) (content string) {
	content, err := d.ReadFileInDirectoryAsString(path...)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return content
}

func (d *DirectoryBase) MustSetParentDirectoryForBaseClass(parentDirectoryForBaseClass Directory) {
	err := d.SetParentDirectoryForBaseClass(parentDirectoryForBaseClass)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (d *DirectoryBase) ReadFileInDirectoryAsLines(path ...string) (content []string, err error) {
	if len(path) == 0 {
		return nil, TracedError("path is empty")
	}

	parent, err := d.GetParentDirectoryForBaseClass()
	if err != nil {
		return nil, err
	}

	fileToRead, err := parent.GetFileInDirectory(path...)
	if err != nil {
		return nil, err
	}

	content, err = fileToRead.ReadAsLines()
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (d *DirectoryBase) ReadFileInDirectoryAsString(path ...string) (content string, err error) {
	parent, err := d.GetParentDirectoryForBaseClass()
	if err != nil {
		return "", err
	}

	fileToRead, err := parent.GetFileInDirectory(path...)
	if err != nil {
		return "", err
	}

	content, err = fileToRead.ReadAsString()
	if err != nil {
		return "", err
	}

	return content, nil
}

func (d *DirectoryBase) SetParentDirectoryForBaseClass(parentDirectoryForBaseClass Directory) (err error) {
	if parentDirectoryForBaseClass == nil {
		return TracedErrorNil("parentDirectoryForBaseClass")
	}

	d.parentDirectoryForBaseClass = parentDirectoryForBaseClass

	return nil
}
