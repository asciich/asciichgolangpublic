# packaga binaryinfo

The aim of this package is to provide information about the build golang binary itself.

## Set software name and version:

If your binary includes this project:

1. You can set a fallback software name which is used in case of no other name is set in the main:
```golang
func main() {
    err := binaryinfo.SetFallbackSoftwareName("a-helper")
	if err != nil {
		panic(err)
	}
}
```
1. Set the version and name using ldflags:
```bash
go build -ldflags="-X 'github.com/asciich/asciichgolangpublic/pkg/binaryinfo.SoftwareName=$(SOFTWARE_NAME)' -X 'github.com/asciich/asciichgolangpublic/pkg/binaryinfo.SoftwareVersion=$(SOFTWARE_VERSION)'" main.go
```
