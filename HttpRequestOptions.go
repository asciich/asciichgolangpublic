package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

// Obsolete: Use https.RequestOptions instead
type HttpRequestOptions struct {
	URL               string
	Verbose           bool
	OutputPath        string
	OverwriteExisting bool
}

// Obsolete: Use https.RequestOptions instead
func NewHttpRequestOptions() (h *HttpRequestOptions) {
	return new(HttpRequestOptions)
}

func (h *HttpRequestOptions) GetOutputPath() (outputPath string, err error) {
	if h.OutputPath == "" {
		return "", tracederrors.TracedErrorf("OutputPath not set")
	}

	return h.OutputPath, nil
}

func (h *HttpRequestOptions) GetOverwriteExisting() (overwriteExisting bool, err error) {

	return h.OverwriteExisting, nil
}

func (h *HttpRequestOptions) GetURL() (uRL string, err error) {
	if h.URL == "" {
		return "", tracederrors.TracedErrorf("URL not set")
	}

	return h.URL, nil
}

func (h *HttpRequestOptions) GetVerbose() (verbose bool, err error) {

	return h.Verbose, nil
}

func (h *HttpRequestOptions) MustGetOutputFile() (outputFile File) {
	outputFile, err := h.GetOutputFile()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return outputFile
}

func (h *HttpRequestOptions) MustGetOutputFilePath() (filePath string) {
	filePath, err := h.GetOutputFilePath()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return filePath
}

func (h *HttpRequestOptions) MustGetOutputPath() (outputPath string) {
	outputPath, err := h.GetOutputPath()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return outputPath
}

func (h *HttpRequestOptions) MustGetOverwriteExisting() (overwriteExisting bool) {
	overwriteExisting, err := h.GetOverwriteExisting()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return overwriteExisting
}

func (h *HttpRequestOptions) MustGetURL() (uRL string) {
	uRL, err := h.GetURL()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return uRL
}

func (h *HttpRequestOptions) MustGetUrl() (url *URL) {
	url, err := h.GetUrl()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return url
}

func (h *HttpRequestOptions) MustGetUrlAsString() (url string) {
	url, err := h.GetUrlAsString()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return url
}

func (h *HttpRequestOptions) MustGetVerbose() (verbose bool) {
	verbose, err := h.GetVerbose()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return verbose
}

func (h *HttpRequestOptions) MustSetOutputPath(outputPath string) {
	err := h.SetOutputPath(outputPath)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (h *HttpRequestOptions) MustSetOutputPathByFile(file File) {
	err := h.SetOutputPathByFile(file)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (h *HttpRequestOptions) MustSetOverwriteExisting(overwriteExisting bool) {
	err := h.SetOverwriteExisting(overwriteExisting)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (h *HttpRequestOptions) MustSetURL(uRL string) {
	err := h.SetURL(uRL)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (h *HttpRequestOptions) MustSetVerbose(verbose bool) {
	err := h.SetVerbose(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (h *HttpRequestOptions) SetOutputPath(outputPath string) (err error) {
	if outputPath == "" {
		return tracederrors.TracedErrorf("outputPath is empty string")
	}

	h.OutputPath = outputPath

	return nil
}

func (h *HttpRequestOptions) SetOverwriteExisting(overwriteExisting bool) (err error) {
	h.OverwriteExisting = overwriteExisting

	return nil
}

func (h *HttpRequestOptions) SetURL(uRL string) (err error) {
	if uRL == "" {
		return tracederrors.TracedErrorf("uRL is empty string")
	}

	h.URL = uRL

	return nil
}

func (h *HttpRequestOptions) SetVerbose(verbose bool) (err error) {
	h.Verbose = verbose

	return nil
}

func (o *HttpRequestOptions) GetDeepCopy() (copy *HttpRequestOptions) {
	copy = NewHttpRequestOptions()
	*copy = *o

	return copy
}

func (o *HttpRequestOptions) GetOutputFile() (outputFile File, err error) {
	filePath, err := o.GetOutputFilePath()
	if err != nil {
		return nil, err
	}

	outputFile, err = GetLocalFileByPath(filePath)
	if err != nil {
		return nil, err
	}

	return outputFile, nil
}

func (o *HttpRequestOptions) GetOutputFilePath() (filePath string, err error) {
	if len(o.OutputPath) > 0 {
		return o.OutputPath, nil
	}

	url, err := o.GetUrl()
	if err != nil {
		return "", err
	}

	filePath, err = url.GetPathBasename()
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func (o *HttpRequestOptions) GetUrl() (url *URL, err error) {
	urlString, err := o.GetUrlAsString()
	if err != nil {
		return nil, err
	}

	url, err = GetUrlFromString(urlString)
	if err != nil {
		return nil, err
	}

	return url, nil
}

func (o *HttpRequestOptions) GetUrlAsString() (url string, err error) {
	if len(o.URL) <= 0 {
		return "", tracederrors.TracedError("Url not set")
	}

	return o.URL, nil
}

func (o *HttpRequestOptions) SetOutputPathByFile(file File) (err error) {
	if file == nil {
		return tracederrors.TracedErrorNil("file")
	}

	localPath, err := file.GetLocalPath()
	if err != nil {
		return err
	}

	err = o.SetOutputPath(localPath)
	if err != nil {
		return err
	}

	return nil
}
