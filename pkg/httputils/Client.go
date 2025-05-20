package httputils

import "github.com/asciich/asciichgolangpublic/files"

type Client interface {
	DownloadAsFile(downloadOptions *DownloadAsFileOptions) (downloadedFile files.File, err error)
	DownloadAsTemporaryFile(downloadOptions *DownloadAsFileOptions) (downloadedFile files.File, err error)
	SendRequest(requestOptions *RequestOptions) (response Response, err error)
	SendRequestAndGetBodyAsString(requestOptions *RequestOptions) (responseBody string, err error)
	SendRequestAndRunYqQueryAgainstBody(requestOptions *RequestOptions, query string) (result string, err error)
}
