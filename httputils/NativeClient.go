package httputils

import (
	"crypto/tls"
	"io"
	"net/http"
	"os"

	"github.com/asciich/asciichgolangpublic/files"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tempfiles"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

// HTTP client written using native go http implementation.
type NativeClient struct {
}

// Get the HTTP client written using native go http implementation.
//
// This is the default client to use when sending request from your running machine.
func GetNativeClient() (client Client) {
	return NewNativeClient()
}

func NewNativeClient() (n *NativeClient) {
	return new(NativeClient)
}

func (c *NativeClient) MustSendRequestAndRunYqQueryAgainstBody(requestOptions *RequestOptions, query string) (result string) {
	result, err := c.SendRequestAndRunYqQueryAgainstBody(requestOptions, query)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return result
}

func (c *NativeClient) SendRequestAndRunYqQueryAgainstBody(requestOptions *RequestOptions, query string) (result string, err error) {
	if requestOptions == nil {
		return "", tracederrors.TracedErrorNil("requestOptions")
	}

	if query == "" {
		return "", tracederrors.TracedErrorEmptyString("query")
	}

	response, err := c.SendRequest(requestOptions)
	if err != nil {
		return "", err
	}

	return response.RunYqQueryAgainstBody(query)
}

func (c *NativeClient) SendRequest(requestOptions *RequestOptions) (response Response, err error) {
	if requestOptions == nil {
		return nil, tracederrors.TracedErrorNil("requestOptions")
	}

	url, err := requestOptions.GetUrl()
	if err != nil {
		return nil, err
	}

	method, err := requestOptions.GetMethodOrDefault()
	if err != nil {
		return nil, err
	}

	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: requestOptions.SkipTLSvalidation}

	client := http.Client{Transport: customTransport}
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	nativeResponse, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer nativeResponse.Body.Close()

	response = NewGenericResponse()
	body, err := io.ReadAll(nativeResponse.Body)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Unable to read body as bytes: %w", err)
	}

	err = response.SetBody(body)
	if err != nil {
		return nil, err
	}

	err = response.SetStatusCode(nativeResponse.StatusCode)
	if err != nil {
		return nil, err
	}

	return response, err
}

func (c *NativeClient) SendRequestAndGetBodyAsString(requestOptions *RequestOptions) (responseBody string, err error) {
	if requestOptions == nil {
		return "", tracederrors.TracedErrorNil("requestOptions")
	}

	response, err := c.SendRequest(requestOptions)
	if err != nil {
		return "", err
	}

	return response.GetBodyAsString()
}

func (n *NativeClient) DownloadAsFile(downloadOptions *DownloadAsFileOptions) (downloadedFile files.File, err error) {
	if downloadOptions == nil {
		return nil, tracederrors.TracedErrorNil("downloadOptions")
	}

	requestOptions, err := downloadOptions.GetRequestOptions()
	if err != nil {
		return nil, err
	}

	url, err := requestOptions.GetUrl()
	if err != nil {
		return nil, err
	}

	request, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer request.Body.Close()

	outputPath, err := downloadOptions.GetOutputPath()
	if err != nil {
		return nil, err
	}

	downloadedFile, err = files.GetLocalFileByPath(outputPath)
	if err != nil {
		return nil, err
	}

	outputFilePath, err := downloadedFile.GetLocalPath()
	if err != nil {
		return nil, err
	}

	if downloadOptions.Sha256Sum != "" {
		exists, err := downloadedFile.Exists(false)
		if err != nil {
			return nil, err
		}

		if exists {
			sha256, err := downloadedFile.GetSha256Sum()
			if err != nil {
				return nil, err
			}

			if sha256 == downloadOptions.Sha256Sum {
				if downloadOptions.Verbose {
					logging.LogInfof(
						"File '%s' already exists and matches sha256sum '%s'. Skip download.",
						outputFilePath,
						sha256,
					)
				}

				return downloadedFile, nil
			}
		}
	}

	if downloadOptions.OverwriteExisting {
		logging.LogInfof("Going to ensure '%s' is absent before download starts", outputFilePath)
		err = downloadedFile.Delete(requestOptions.Verbose)
		if err != nil {
			return nil, err
		}
	}

	if downloadOptions.Verbose {
		logging.LogInfof("Going to download: '%s' as file '%s'.", url, outputFilePath)
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

	if downloadOptions.Sha256Sum != "" {
		expectedSha256 := downloadOptions.Sha256Sum

		if downloadOptions.Verbose {
			logging.LogInfof(
				"Going to validate downloaded file '%s' using expected sha256sum %s",
				outputFilePath,
				expectedSha256,
			)
		}

		sha256, err := downloadedFile.GetSha256Sum()
		if err != nil {
			return nil, err
		}

		if expectedSha256 == sha256 {
			if downloadOptions.Verbose {
				logging.LogInfof(
					"Downloaded file '%s' matches expected sha256sum %s",
					outputFilePath,
					expectedSha256,
				)
			}
		} else {
			return nil, tracederrors.TracedErrorf(
				"Downloaded file '%s' has checksum '%s' and is not matching expected '%s'.",
				outputFilePath,
				sha256,
				expectedSha256,
			)
		}
	}

	return downloadedFile, nil
}

func (n *NativeClient) DownloadAsTemporaryFile(downloadOptions *DownloadAsFileOptions) (downloadedFile files.File, err error) {
	if downloadOptions == nil {
		return nil, tracederrors.TracedErrorNil("downloadOptions")
	}

	toUse := downloadOptions.GetDeepCopy()

	toUse.OutputPath, err = tempfiles.CreateEmptyTemporaryFileAndGetPath(false)
	if err != nil {
		return nil, err
	}
	toUse.OverwriteExisting = true

	return n.DownloadAsFile(toUse)
}

func (n *NativeClient) MustDownloadAsFile(downloadOptions *DownloadAsFileOptions) (downloadedFile files.File) {
	downloadedFile, err := n.DownloadAsFile(downloadOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return downloadedFile
}

func (n *NativeClient) MustDownloadAsTemporaryFile(downloadOptions *DownloadAsFileOptions) (downloadedFile files.File) {
	downloadedFile, err := n.DownloadAsTemporaryFile(downloadOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return downloadedFile
}

func (n *NativeClient) MustSendRequest(requestOptions *RequestOptions) (response Response) {
	response, err := n.SendRequest(requestOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return response
}

func (n *NativeClient) MustSendRequestAndGetBodyAsString(requestOptions *RequestOptions) (responseBody string) {
	responseBody, err := n.SendRequestAndGetBodyAsString(requestOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return responseBody
}
