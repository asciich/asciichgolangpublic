package httpoptions

import (
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type DownloadAsTemporaryFileOptions struct {
	RequestOptions *RequestOptions

	// If Sha256Sum is set:
	// - The download will be skipped if OutputPath has already the expected content.
	// - The download is validated.
	Sha256Sum string
}

func (d *DownloadAsTemporaryFileOptions) GetDeepCopy() (copy *DownloadAsTemporaryFileOptions) {
	copy = new(DownloadAsTemporaryFileOptions)
	*copy = *d

	if d.RequestOptions != nil {
		copy.RequestOptions = d.RequestOptions.GetDeepCopy()
	}

	return copy
}

func (d *DownloadAsTemporaryFileOptions) GetRequestOptions() (requestOptions *RequestOptions, err error) {
	if d.RequestOptions == nil {
		return nil, tracederrors.TracedErrorf("RequestOptions not set")
	}

	return d.RequestOptions, nil
}

func (d *DownloadAsTemporaryFileOptions) SetRequestOptions(requestOptions *RequestOptions) (err error) {
	if requestOptions == nil {
		return tracederrors.TracedErrorf("requestOptions is nil")
	}

	d.RequestOptions = requestOptions

	return nil
}
