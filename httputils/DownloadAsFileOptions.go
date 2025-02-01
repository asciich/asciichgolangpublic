package httputils

import (
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type DownloadAsFileOptions struct {
	RequestOptions    *RequestOptions
	OutputPath        string
	OverwriteExisting bool
	Verbose           bool
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

func (d *DownloadAsFileOptions) GetVerbose() (verbose bool) {

	return d.Verbose
}

func (d *DownloadAsFileOptions) MustGetOutputPath() (outputPath string) {
	outputPath, err := d.GetOutputPath()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return outputPath
}

func (d *DownloadAsFileOptions) MustGetRequestOptions() (requestOptions *RequestOptions) {
	requestOptions, err := d.GetRequestOptions()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return requestOptions
}

func (d *DownloadAsFileOptions) MustSetOutputPath(outputPath string) {
	err := d.SetOutputPath(outputPath)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (d *DownloadAsFileOptions) MustSetRequestOptions(requestOptions *RequestOptions) {
	err := d.SetRequestOptions(requestOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
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

func (d *DownloadAsFileOptions) SetVerbose(verbose bool) {
	d.Verbose = verbose
}
