package httpnativeclientoo

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
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
		transportToUse = http.DefaultTransport.(*http.Transport).Clone()
		transportToUse.TLSClientConfig = &tls.Config{InsecureSkipVerify: requestOptions.SkipTLSvalidation}
	}

	if transportToUse == nil {
		return nil, tracederrors.TracedError("TransportToUse is nil after evaluation.")
	}

	// Store certificates if requested
	var collectedCerts []*x509.Certificate
	if requestOptions.CollectCertificates {
		originalVerifyConnection := transportToUse.TLSClientConfig.VerifyConnection
		transportToUse.TLSClientConfig.VerifyConnection = func(state tls.ConnectionState) error {
			collectedCerts = make([]*x509.Certificate, len(state.PeerCertificates))
			copy(collectedCerts, state.PeerCertificates)
			if originalVerifyConnection != nil {
				return originalVerifyConnection(state)
			}
			return nil
		}
	}

	client := http.Client{Transport: transportToUse}

	logging.LogInfoByCtxf(ctx, "http native client is used to send request to %s", url)

	if requestOptions.Data != nil {
		logging.LogInfoByCtxf(ctx, "The request body of '%d' bytes was added for %s .", len(requestOptions.Data), url)
	}

	// Retry logic for transient server errors (500, 502, 503, 504)
	var nativeResponse *http.Response
	var lastErr error
	maxRetries := 3
	baseDelay := 500 * time.Millisecond

	for attempt := 0; attempt <= maxRetries; attempt++ {
		request, reqErr := http.NewRequest(method, url, nil)
		if reqErr != nil {
			return nil, reqErr
		}

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
		}

		nativeResponse, lastErr = client.Do(request)
		if lastErr != nil {
			// Network error, retry
			if attempt < maxRetries {
				delay := baseDelay * time.Duration(attempt+1)
				logging.LogInfoByCtxf(ctx, "Request failed with error '%v', retrying in %v (attempt %d/%d)...", lastErr, delay, attempt+1, maxRetries+1)
				select {
				case <-ctx.Done():
					return nil, tracederrors.TracedErrorf("Context cancelled during retry: %w", ctx.Err())
				case <-time.After(delay):
				}
				continue
			}
			return nil, lastErr
		}

		// Check if status code is a retryable server error
		if nativeResponse.StatusCode >= 500 && nativeResponse.StatusCode <= 599 {
			if attempt < maxRetries {
				nativeResponse.Body.Close()
				delay := baseDelay * time.Duration(attempt+1)
				logging.LogInfoByCtxf(ctx, "Received server error status code %d, retrying in %v (attempt %d/%d)...", nativeResponse.StatusCode, delay, attempt+1, maxRetries+1)
				select {
				case <-ctx.Done():
					return nil, tracederrors.TracedErrorf("Context cancelled during retry: %w", ctx.Err())
				case <-time.After(delay):
				}
				continue
			}
		}

		// Success or non-retryable error
		break
	}

	defer nativeResponse.Body.Close()

	response = httpgeneric.NewGenericResponse()
	body, readErr := io.ReadAll(nativeResponse.Body)
	if readErr != nil {
		return nil, tracederrors.TracedErrorf("Unable to read body as bytes: %w", readErr)
	}

	err = response.SetBody(body)
	if err != nil {
		return nil, err
	}

	err = response.SetStatusCode(nativeResponse.StatusCode)
	if err != nil {
		return nil, err
	}

	// Set collected certificates if requested
	if requestOptions.CollectCertificates && collectedCerts != nil {
		if genericResp, ok := response.(*httpgeneric.GenericResponse); ok {
			err = genericResp.SetServerCertificates(collectedCerts)
			if err != nil {
				return nil, err
			}
		}
	}

	err = response.CheckStatusCode([]int{http.StatusOK, http.StatusCreated})
	if err != nil {
		return response, err
	}

	return response, nil
}

func (c *NativeClient) SendRequestAndGetBodyAsBytes(ctx context.Context, requestOptions *httpoptions.RequestOptions) (responseBody []byte, err error) {
	if requestOptions == nil {
		return nil, tracederrors.TracedErrorNil("requestOptions")
	}

	response, err := c.SendRequest(ctx, requestOptions)
	if err != nil {
		return nil, err
	}

	return response.GetBodyAsBytes()
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

	targetFilePath, err := downloadedFile.GetLocalPath()
	if err != nil {
		return nil, err
	}

	if downloadOptions.Sha256Sum != "" {
		exists, err := downloadedFile.Exists(contextutils.WithSilent(ctx))
		if err != nil {
			return nil, err
		}

		if exists {
			sha256, err := downloadedFile.GetSha256Sum(ctx)
			if err != nil {
				return nil, err
			}

			if sha256 == downloadOptions.Sha256Sum {
				logging.LogInfoByCtxf(ctx, "File '%s' already exists and matches sha256sum '%s'. Skip download.", targetFilePath, sha256)

				return downloadedFile, nil
			}
		}
	}

	var outputFilePath string
	if downloadOptions.UseSudo {
		outputFilePath, err = tempfiles.CreateTemporaryFile(ctx)
		if err != nil {
			return nil, err
		}
	} else {
		outputFilePath, err = downloadedFile.GetLocalPath()
		if err != nil {
			return nil, err
		}
	}

	if downloadOptions.OverwriteExisting {
		logging.LogInfoByCtxf(ctx, "Going to ensure '%s' is absent before download starts", outputFilePath)
		err = downloadedFile.Delete(ctx, &filesoptions.DeleteOptions{
			UseSudo: downloadOptions.UseSudo,
		})
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

	if downloadOptions.UseSudo {
		err = nativefiles.Copy(ctx, outputFilePath, targetFilePath, &filesoptions.CopyOptions{UseSudo: downloadOptions.UseSudo})
		if err != nil {
			return nil, err
		}
	}

	if downloadOptions.PermissionsString != "" {
		err = nativefiles.Chmod(ctx, targetFilePath, &filesoptions.ChmodOptions{
			PermissionsString: downloadOptions.PermissionsString,
			UseSudo:           downloadOptions.UseSudo,
		})
		if err != nil {
			return nil, err
		}
	}

	logging.LogChangedByCtxf(ctx, "Downloaded '%s' as file '%s'.", url, targetFilePath)

	if downloadOptions.Sha256Sum != "" {
		expectedSha256 := downloadOptions.Sha256Sum

		logging.LogInfoByCtxf(ctx, "Going to validate downloaded file '%s' using expected sha256sum %s", targetFilePath, expectedSha256)

		sha256, err := downloadedFile.GetSha256Sum(ctx)
		if err != nil {
			return nil, err
		}

		if expectedSha256 == sha256 {
			logging.LogInfoByCtxf(ctx, "Downloaded file '%s' matches expected sha256sum %s", targetFilePath, expectedSha256)
		} else {
			return nil, tracederrors.TracedErrorf(
				"%w: Downloaded file '%s' has checksum '%s' and is not matching expected '%s'.",
				httpgeneric.ErrChecksumMismatch,
				targetFilePath,
				sha256,
				expectedSha256,
			)
		}
	}

	return downloadedFile, nil
}

func (n *NativeClient) DownloadAsTemporaryFile(ctx context.Context, downloadOptions *httpoptions.DownloadAsTemporaryFileOptions) (downloadedFile filesinterfaces.File, err error) {
	if downloadOptions == nil {
		return nil, tracederrors.TracedErrorNil("downloadOptions")
	}

	outputPath, err := tempfiles.CreateTemporaryFile(contextutils.WithSilent(ctx))
	if err != nil {
		return nil, err
	}

	options := &httpoptions.DownloadAsFileOptions{
		RequestOptions:    downloadOptions.RequestOptions,
		Sha256Sum:         downloadOptions.Sha256Sum,
		OutputPath:        outputPath,
		OverwriteExisting: true,
	}

	return n.DownloadAsFile(ctx, options)
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
