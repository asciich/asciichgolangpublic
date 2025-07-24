package httputils

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"os"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfilesoo"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httputilsimplementationindependend"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httputilsinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httputilsparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// HTTP client written using native go http implementation.
type NativeClient struct {
}

// Get the HTTP client written using native go http implementation.
//
// This is the default client to use when sending request from your running machine.
func GetNativeClient() (client httputilsinterfaces.Client) {
	return NewNativeClient()
}

func NewNativeClient() (n *NativeClient) {
	return new(NativeClient)
}

func (c *NativeClient) SendRequestAndRunYqQueryAgainstBody(ctx context.Context, requestOptions *httputilsparameteroptions.RequestOptions, query string) (result string, err error) {
	if requestOptions == nil {
		return "", tracederrors.TracedErrorNil("requestOptions")
	}

	if query == "" {
		return "", tracederrors.TracedErrorEmptyString("query")
	}

	response, err := c.SendRequest(ctx, requestOptions)
	if err != nil {
		return "", err
	}

	return response.RunYqQueryAgainstBody(query)
}

func (c *NativeClient) SendRequest(ctx context.Context, requestOptions *httputilsparameteroptions.RequestOptions) (response httputilsinterfaces.Response, err error) {
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

	response = httputilsimplementationindependend.NewGenericResponse()
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

	err = response.CheckStatusCode(httputilsimplementationindependend.STATUS_CODE_OK)
	if err != nil {
		return response, err
	}

	return response, err
}

func (c *NativeClient) SendRequestAndGetBodyAsString(ctx context.Context, requestOptions *httputilsparameteroptions.RequestOptions) (responseBody string, err error) {
	if requestOptions == nil {
		return "", tracederrors.TracedErrorNil("requestOptions")
	}

	response, err := c.SendRequest(ctx, requestOptions)
	if err != nil {
		return "", err
	}

	return response.GetBodyAsString()
}

func (n *NativeClient) DownloadAsFile(ctx context.Context, downloadOptions *httputilsparameteroptions.DownloadAsFileOptions) (downloadedFile files.File, err error) {
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
				logging.LogInfoByCtxf(ctx, "File '%s' already exists and matches sha256sum '%s'. Skip download.", outputFilePath, sha256)

				return downloadedFile, nil
			}
		}
	}

	if downloadOptions.OverwriteExisting {
		logging.LogInfoByCtxf(ctx, "Going to ensure '%s' is absent before download starts", outputFilePath)
		err = downloadedFile.Delete(contextutils.GetVerboseFromContext(ctx))
		if err != nil {
			return nil, err
		}
	}

	logging.LogInfoByCtxf(ctx, "Going to download: '%s' as file '%s'.", url, outputFilePath)

	outFd, err := os.Create(outputFilePath)
	if err != nil {
		return nil, tracederrors.TracedError(err.Error())
	}
	defer outFd.Close()

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	chunkSize := GetProgressEveryNBytes(ctx)
	if chunkSize <= 0 {
		outFd.ReadFrom(response.Body)
	} else {
		buf := make([]byte, chunkSize)
		var downloadedBytes int64
		var totalBytes = response.ContentLength
		var eofDetected bool
		for {
			n, err := response.Body.Read(buf)
			if err != nil {
				if err == io.EOF {
					eofDetected = true
				} else {
					return nil, tracederrors.TracedErrorf("Error while downloading: %w", err)
				}
			}
			if n > 0 {
				_, err = outFd.Write(buf[:n])
				if err != nil {
					return nil, tracederrors.TracedErrorf("Error while writing downloaded data to file '%s': %w", outputFilePath, err)
				}
				downloadedBytes += int64(n)
				progressPercent := 100. / float64(totalBytes) * float64(downloadedBytes)
				logging.LogInfoByCtxf(ctx, "Downloaded %d/%d bytes (%.02f%%)", downloadedBytes, totalBytes, progressPercent)
			}

			if eofDetected {
				break
			}
		}
	}

	logging.LogInfoByCtxf(ctx, "Downloaded '%s' as file '%s'.", url, outputFilePath)

	if downloadOptions.Sha256Sum != "" {
		expectedSha256 := downloadOptions.Sha256Sum

		logging.LogInfoByCtxf(ctx, "Going to validate downloaded file '%s' using expected sha256sum %s", outputFilePath, expectedSha256)

		sha256, err := downloadedFile.GetSha256Sum()
		if err != nil {
			return nil, err
		}

		if expectedSha256 == sha256 {
			logging.LogInfoByCtxf(ctx, "Downloaded file '%s' matches expected sha256sum %s", outputFilePath, expectedSha256)
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

func (n *NativeClient) DownloadAsTemporaryFile(ctx context.Context, downloadOptions *httputilsparameteroptions.DownloadAsFileOptions) (downloadedFile files.File, err error) {
	if downloadOptions == nil {
		return nil, tracederrors.TracedErrorNil("downloadOptions")
	}

	toUse := downloadOptions.GetDeepCopy()

	toUse.OutputPath, err = tempfilesoo.CreateEmptyTemporaryFileAndGetPath(false)
	if err != nil {
		return nil, err
	}
	toUse.OverwriteExisting = true

	return n.DownloadAsFile(ctx, toUse)
}
