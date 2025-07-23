package httputilsparameteroptions

import (
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/tracederrors"
)

type DownloadAsFileOptions struct {
	RequestOptions    *RequestOptions
	OutputPath        string
	OverwriteExisting bool

	// If Sha256Sum is set:
	// - The download will be skipped if OutputPath has already the expected content.
	// - The download is validated.
	Sha256Sum string
}

func NewDownloadAsFileOptions() (d *DownloadAsFileOptions) {
	return new(DownloadAsFileOptions)
}

func (d *DownloadAsFileOptions) GetDeepCopy() (copy *DownloadAsFileOptions) {
	copy = new(DownloadAsFileOptions)
	*copy = *d

	if d.RequestOptions != nil {
		copy.RequestOptions = d.RequestOptions.GetDeepCopy()
	}

	return copy
}

func (d *DownloadAsFileOptions) GetOutputPath() (outputPath string, err error) {
	if d.OutputPath == "" {
		return "", tracederrors.TracedErrorf("OutputPath not set")
	}

	return d.OutputPath, nil
}

func (d *DownloadAsFileOptions) GetOverwriteExisting() (overwriteExisting bool) {

	return d.OverwriteExisting
}

func (d *DownloadAsFileOptions) GetRequestOptions() (requestOptions *RequestOptions, err error) {
	if d.RequestOptions == nil {
		return nil, tracederrors.TracedErrorf("RequestOptions not set")
	}

	return d.RequestOptions, nil
}

func (d *DownloadAsFileOptions) SetOutputPath(outputPath string) (err error) {
	if outputPath == "" {
		return tracederrors.TracedErrorf("outputPath is empty string")
	}

	d.OutputPath = outputPath

	return nil
}

func (d *DownloadAsFileOptions) SetOverwriteExisting(overwriteExisting bool) {
	d.OverwriteExisting = overwriteExisting
}

func (d *DownloadAsFileOptions) SetRequestOptions(requestOptions *RequestOptions) (err error) {
	if requestOptions == nil {
		return tracederrors.TracedErrorf("requestOptions is nil")
	}

	d.RequestOptions = requestOptions

	return nil
}
