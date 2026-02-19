package httpnativeclientoo

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"os"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfilesoo"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpoptions"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httputilsinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"github.com/asciich/asciichgolangpublic/pkg/urlsutils"
)

// HTTP client written using native go http implementation.
type NativeClient struct {
	baseUrl   string
	port      int
	basicAuth *httpoptions.BasicAuth
}

// Get the HTTP client written using native go http implementation.
//
// This is the default client to use when sending request from your running machine.
func NewNativeClient() (n *NativeClient) {
	return new(NativeClient)
}

func (c *NativeClient) SendRequestAndRunYqQueryAgainstBody(ctx context.Context, requestOptions *httpoptions.RequestOptions, query string) (result string, err error) {
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

func (c *NativeClient) SendRequest(ctx context.Context, requestOptions *httpoptions.RequestOptions) (response httputilsinterfaces.Response, err error) {
	if requestOptions == nil {
		return nil, tracederrors.TracedErrorNil("requestOptions")
	}

	url := c.baseUrl

	if url == "" {
		url, err = requestOptions.GetUrl()
		if err != nil {
			return nil, err
		}
	}

	if requestOptions.Path != "" {
		url, err = urlsutils.SetPath(url, requestOptions.Path)
		if err != nil {
			return nil, err
		}
	}

	if url == "" {
		return nil, tracederrors.TracedError("url is empty string after evaluation")
	}

	if requestOptions.Port != 0 {
		url, err = urlsutils.SetPort(url, requestOptions.Port)
		if err != nil {
			return nil, err
		}
	} else {
		if c.port != 0 {
			url, err = urlsutils.SetPort(url, c.port)
			if err != nil {
				return nil, err
			}
		}
	}

	method, err := requestOptions.GetMethodOrDefault()
	if err != nil {
		return nil, err
	}

	var transportToUse *http.Transport
	if requestOptions.TransportToUse != nil {
		transportToUse = requestOptions.TransportToUse
	} else {
		transportToUse := http.DefaultTransport.(*http.Transport).Clone()
		transportToUse.TLSClientConfig = &tls.Config{InsecureSkipVerify: requestOptions.SkipTLSvalidation}
	}

	client := http.Client{Transport: transportToUse}

	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "http native client is used to send request to %s", request.URL.String())

	if requestOptions.BasicAuth != nil {
		request.Header.Set("Authorization", requestOptions.BasicAuth.AuthorizationValue())
	} else {
		if c.basicAuth != nil {
			request.Header.Set("Authorization", c.basicAuth.AuthorizationValue())
		}
	}

	for k, v := range requestOptions.Header {
		request.Header.Set(k, v)
	}

	if requestOptions.Data != nil {
		request.Body = io.NopCloser(bytes.NewReader(requestOptions.Data))
		request.ContentLength = int64(len(requestOptions.Data))
		logging.LogInfoByCtxf(ctx, "The request body of '%d' bytes was added for %s .", request.ContentLength, url)
	}

	nativeResponse, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer nativeResponse.Body.Close()

	response = httpgeneric.NewGenericResponse()
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

	err = response.CheckStatusCode(http.StatusOK)
	if err != nil {
		return response, err
	}

	return response, err
}

func (c *NativeClient) SendRequestAndGetBodyAsString(ctx context.Context, requestOptions *httpoptions.RequestOptions) (responseBody string, err error) {
	if requestOptions == nil {
		return "", tracederrors.TracedErrorNil("requestOptions")
	}

	response, err := c.SendRequest(ctx, requestOptions)
	if err != nil {
		return "", err
	}

	return response.GetBodyAsString()
}

func (n *NativeClient) DownloadAsFile(ctx context.Context, downloadOptions *httpoptions.DownloadAsFileOptions) (downloadedFile filesinterfaces.File, err error) {
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
		exists, err := downloadedFile.Exists(contextutils.ContextSilent())
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
		err = downloadedFile.Delete(ctx, &filesoptions.DeleteOptions{})
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

	chunkSize := httpgeneric.GetProgressEveryNBytes(ctx)
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

func (n *NativeClient) DownloadAsTemporaryFile(ctx context.Context, downloadOptions *httpoptions.DownloadAsFileOptions) (downloadedFile filesinterfaces.File, err error) {
	if downloadOptions == nil {
		return nil, tracederrors.TracedErrorNil("downloadOptions")
	}

	toUse := downloadOptions.GetDeepCopy()

	toUse.OutputPath, err = tempfilesoo.CreateEmptyTemporaryFileAndGetPath(contextutils.WithSilent(ctx))
	if err != nil {
		return nil, err
	}
	toUse.OverwriteExisting = true

	return n.DownloadAsFile(ctx, toUse)
}

func (n *NativeClient) SetPort(port int) error {
	if port <= 0 {
		return tracederrors.TracedErrorf("Invalid port '%d'", port)
	}

	n.port = port

	return nil
}

func (n *NativeClient) SetBasicAuth(basicAuth *httpoptions.BasicAuth) error {
	if basicAuth == nil {
		return tracederrors.TracedErrorNil("basicAuth")
	}

	n.basicAuth = basicAuth

	return nil
}

func (n *NativeClient) SetBaseUrl(baseUrl string) error {
	err := urlsutils.CheckIsUrl(baseUrl)
	if err != nil {
		return err
	}

	n.baseUrl = baseUrl

	return nil
}
