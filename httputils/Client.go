package httputils

import "github.com/asciich/asciichgolangpublic/files"

type Client interface {
	DownloadAsFile(downloadOptions *DownloadAsFileOptions) (downloadedFile files.File, err error)
	DownloadAsTemporaryFile(downloadOptions *DownloadAsFileOptions) (downloadedFile files.File, err error)
	MustDownloadAsFile(downloadOptions *DownloadAsFileOptions) (downloadedFile files.File)
	MustDownloadAsTemporaryFile(downloadOptions *DownloadAsFileOptions) (downloadedFile files.File)
	MustSendRequest(requestOptions *RequestOptions) (response Response)
	MustSendRequestAndGetBodyAsString(requestOptions *RequestOptions) (responseBody string)
	SendRequest(requestOptions *RequestOptions) (response Response, err error)
	SendRequestAndGetBodyAsString(requestOptions *RequestOptions) (responseBody string, err error)
}
