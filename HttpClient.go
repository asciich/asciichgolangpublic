package asciichgolangpublic

import (
	"net/http"
	"os"

	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type HttpClientService struct{}

// Obsolete: use http.NativeClient() instead:
func HttpClient() (httpClient *HttpClientService) {
	return new(HttpClientService)
}

// Obsolete: use http.NativeClient() instead:
func NewHttpClientService() (h *HttpClientService) {
	return new(HttpClientService)
}

func (h *HttpClientService) DownloadAsFile(requestOptions *HttpRequestOptions) (downloadedFile File, err error) {
	if requestOptions == nil {
		return nil, tracederrors.TracedErrorNil("requestOptions")
	}

	url, err := requestOptions.GetUrlAsString()
	if err != nil {
		return nil, err
	}

	request, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer request.Body.Close()

	downloadedFile, err = requestOptions.GetOutputFile()
	if err != nil {
		return nil, err
	}

	outputFilePath, err := downloadedFile.GetLocalPath()
	if err != nil {
		return nil, err
	}

	if requestOptions.OverwriteExisting {
		logging.LogInfof("Going to ensure '%s' is absent before download starts", outputFilePath)
		err = downloadedFile.Delete(requestOptions.Verbose)
		if err != nil {
			return nil, err
		}
	}

	if requestOptions.Verbose {
		logging.LogInfof("Going to download: '%v' as file '%v'.", url, outputFilePath)
	}

	outFd, err := os.Create(outputFilePath)
	if err != nil {
		return nil, tracederrors.TracedError(err.Error())
	}
	defer outFd.Close()
	outFd.ReadFrom(request.Body)

	if requestOptions.Verbose {
		logging.LogInfof("Downloaded '%v' as file '%v'.", url, outputFilePath)
	}

	return downloadedFile, nil
}

func (h *HttpClientService) DownloadAsTemporaryFile(requestOptions *HttpRequestOptions) (downloadedFile File, err error) {
	if requestOptions == nil {
		return nil, tracederrors.TracedErrorNil("requestOptions")
	}

	requestOptionsToUse := requestOptions.GetDeepCopy()

	temporaryFile, err := TemporaryFiles().CreateEmptyTemporaryFile(requestOptions.Verbose)
	if err != nil {
		return nil, err
	}

	err = requestOptionsToUse.SetOverwriteExisting(true)
	if err != nil {
		return nil, err
	}

	err = requestOptionsToUse.SetOutputPathByFile(temporaryFile)
	if err != nil {
		return nil, err
	}

	downloadedFile, err = h.DownloadAsFile(requestOptionsToUse)
	if err != nil {
		return nil, err
	}

	return downloadedFile, nil
}

func (h *HttpClientService) MustDownloadAsFile(requestOptions *HttpRequestOptions) (downloadedFile File) {
	downloadedFile, err := h.DownloadAsFile(requestOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return downloadedFile
}

func (h *HttpClientService) MustDownloadAsTemporaryFile(requestOptions *HttpRequestOptions) (downloadedFile File) {
	downloadedFile, err := h.DownloadAsTemporaryFile(requestOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return downloadedFile
}
