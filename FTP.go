package asciichgolangpublic

type FTPService struct{}

func FTP() (f *FTPService) {
	return NewFTPService()
}

func NewFTPService() (f *FTPService) {
	return new(FTPService)
}

func (f *FTPService) GetDefaultPort() (defaultPort int) {
	return 21
}
