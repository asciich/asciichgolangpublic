package http

import "github.com/asciich/asciichgolangpublic/files"

type Client interface {
	DownloadAsFile(downloadOptions *DownloadAsFileOptions) (downloadedFile files.File, err error)
	DownloadAsTemporaryFile(downloadOptions *DownloadAsFileOptions) (downloadedFile files.File, err error)
	SendRequest(requestOptions *RequestOptions) (response Response, err error)
	MustDownloadAsFile(downloadOptions *DownloadAsFileOptions) (downloadedFile files.File)
	MustDownloadAsTemporaryFile(downloadOptions *DownloadAsFileOptions) (downloadedFile files.File)
	MustSendRequest(requestOptions *RequestOptions) (response Response)
}
