package httpcommandexecutorclientoo

import (
	"context"
	"net/http"

	"github.com/asciich/asciichgolangpublic/pkg/applications/curlutils"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfileoo"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutortempfile"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpoptions"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httputilsinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"github.com/asciich/asciichgolangpublic/pkg/urlsutils"
)

type HttpCommandExecutorClient struct {
	commandExecutor commandexecutorinterfaces.CommandExecutor
	port            int
}

func NewClient(commandExecutor commandexecutorinterfaces.CommandExecutor) (*HttpCommandExecutorClient, error) {
	if commandExecutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExecutor")
	}

	client := &HttpCommandExecutorClient{
		commandExecutor: commandExecutor,
	}

	return client, nil
}

func (h *HttpCommandExecutorClient) GetCommandExecutor() (commandexecutorinterfaces.CommandExecutor, error) {
	if h.commandExecutor == nil {
		return nil, tracederrors.TracedError("commandExecutor not set")
	}

	return h.commandExecutor, nil
}

func (h *HttpCommandExecutorClient) DownloadAsFile(ctx context.Context, downloadOptions *httpoptions.DownloadAsFileOptions) (downloadedFile filesinterfaces.File, err error) {
	if downloadOptions == nil {
		return nil, tracederrors.TracedErrorNil("downloadOptions")
	}

	commandExecutor, err := h.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return nil, err
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

	downloadedFile, err = commandexecutorfileoo.New(commandExecutor, outputPath)
	if err != nil {
		return nil, err
	}

	outputFilePath, err := downloadedFile.GetPath()
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
				logging.LogInfoByCtxf(ctx, "File '%s' already exists and matches sha256sum '%s'. Skip download.", outputFilePath, sha256)

				return downloadedFile, nil
			}
		}
	}

	if downloadOptions.OverwriteExisting {
		logging.LogInfoByCtxf(ctx, "Going to ensure '%s' is absent before download starts", outputFilePath)
		err = downloadedFile.Delete(ctx, &filesoptions.DeleteOptions{UseSudo: downloadOptions.UseSudo})
		if err != nil {
			return nil, err
		}
	}

	logging.LogInfoByCtxf(ctx, "Going to download: '%s' as file '%s' on '%s'.", url, outputFilePath, hostDescription)

	command := []string{"curl", "--fail"}

	if downloadOptions.RequestOptions.SkipTLSvalidation {
		command = append(command, "--insecure")
	}

	command = append(command, "-L", url, "--output", outputFilePath)

	if downloadOptions.UseSudo {
		command = append([]string{"sudo"}, command...)
	}

	_, err = commandExecutor.RunCommand(ctx, &parameteroptions.RunCommandOptions{
		Command: command,
	})
	if err != nil {
		return nil, err
	}

	logging.LogChangedByCtxf(ctx, "Downloaded '%s' as file '%s' on '%s'.", url, outputFilePath, hostDescription)

	if downloadOptions.Sha256Sum != "" {
		expectedSha256 := downloadOptions.Sha256Sum

		logging.LogInfoByCtxf(ctx, "Going to validate downloaded file '%s' using expected sha256sum %s", outputFilePath, expectedSha256)

		sha256, err := downloadedFile.GetSha256Sum(ctx)
		if err != nil {
			return nil, err
		}

		if expectedSha256 == sha256 {
			logging.LogInfoByCtxf(ctx, "Downloaded file '%s' matches expected sha256sum %s", outputFilePath, expectedSha256)
		} else {
			return nil, tracederrors.TracedErrorf(
				"%w: Downloaded file '%s' has checksum '%s' and is not matching expected '%s'.",
				httpgeneric.ErrChecksumMismatch,
				outputFilePath,
				sha256,
				expectedSha256,
			)
		}
	}

	return downloadedFile, nil
}

func (h *HttpCommandExecutorClient) DownloadAsTemporaryFile(ctx context.Context, downloadOptions *httpoptions.DownloadAsTemporaryFileOptions) (downloadedFile filesinterfaces.File, err error) {
	if downloadOptions == nil {
		return nil, tracederrors.TracedErrorNil("downloadOptions")
	}

	commandExecutor, err := h.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	outputPath, err := commandexecutortempfile.CreateEmptyTemporaryFile(ctx, commandExecutor)
	if err != nil {
		return nil, err
	}

	options := &httpoptions.DownloadAsFileOptions{
		RequestOptions:    downloadOptions.RequestOptions,
		Sha256Sum:         downloadOptions.Sha256Sum,
		OutputPath:        outputPath,
		OverwriteExisting: true,
	}

	return h.DownloadAsFile(ctx, options)
}

func (h *HttpCommandExecutorClient) SendRequest(ctx context.Context, requestOptions *httpoptions.RequestOptions) (response httputilsinterfaces.Response, err error) {
	if requestOptions == nil {
		return nil, tracederrors.TracedErrorNil("requestOptions")
	}

	commandExecutor, err := h.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	headersFile, err := commandexecutortempfile.CreateEmptyTemporaryFile(ctx, commandExecutor)
	if err != nil {
		return nil, err
	}
	defer commandexecutorfile.Delete(ctx, commandExecutor, headersFile, &filesoptions.DeleteOptions{})

	method, err := requestOptions.GetMethodOrDefault()
	if err != nil {
		return nil, err
	}

	url, err := requestOptions.GetUrl()
	if err != nil {
		return nil, err
	}

	if requestOptions.Port != 0 {
		url, err = urlsutils.SetPort(url, requestOptions.Port)
		if err != nil {
			return nil, err
		}
	} else {
		if h.port != 0 {
			url, err = urlsutils.SetPort(url, h.port)
			if err != nil {
				return nil, err
			}
		}
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Send '%s' request to '%s' on '%s' started.", method, url, hostDescription)

	command := []string{"curl", "--dump-header", headersFile, "-X", method}

	if requestOptions.SkipTLSvalidation {
		command = append(command, "--insecure")
	}

	command = append(command, url)
	output, err := commandExecutor.RunCommand(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: command,
		},
	)
	if err != nil {
		return nil, err
	}

	stdout, err := output.GetStdoutAsBytes()
	if err != nil {
		return nil, err
	}

	response = &httpgeneric.GenericResponse{}
	err = response.SetBody(stdout)
	if err != nil {
		return nil, err
	}

	headersContent, err := commandexecutorfile.ReadAsBytes(commandExecutor, headersFile)
	if err != nil {
		return nil, err
	}

	responseHeader, err := curlutils.ParseHttpHeader(headersContent)
	if err != nil {
		return nil, err
	}

	statusCode, err := responseHeader.GetStatusCode()
	if err != nil {
		return nil, err
	}

	err = response.SetStatusCode(statusCode)
	if err != nil {
		return nil, err
	}

	err = response.CheckStatusCode([]int{http.StatusOK, http.StatusCreated})
	if err != nil {
		return response, err
	}

	logging.LogInfoByCtxf(ctx, "Send '%s' request to '%s' on '%s' started.", method, url, hostDescription)

	return response, nil
}

func (h *HttpCommandExecutorClient) SendRequestAndGetBodyAsString(ctx context.Context, requestOptions *httpoptions.RequestOptions) (responseBody string, err error) {
	response, err := h.SendRequest(ctx, requestOptions)
	if err != nil {
		return "", err
	}

	return response.GetBodyAsString()
}
func (h *HttpCommandExecutorClient) SendRequestAndRunYqQueryAgainstBody(ctx context.Context, requestOptions *httpoptions.RequestOptions, query string) (result string, err error) {
	if requestOptions == nil {
		return "", tracederrors.TracedErrorNil("requestOptions")
	}

	if query == "" {
		return "", tracederrors.TracedErrorEmptyString("query")
	}

	response, err := h.SendRequest(ctx, requestOptions)
	if err != nil {
		return "", err
	}

	return response.RunYqQueryAgainstBody(query)
}
func (h *HttpCommandExecutorClient) SetBaseUrl(baseUrl string) error {
	return tracederrors.TracedErrorNotImplemented()
}
func (h *HttpCommandExecutorClient) SetPort(port int) error {
	if port <= 0 {
		return tracederrors.TracedErrorf("Invalid port: %d", port)
	}

	h.port = port

	return nil
}
func (h *HttpCommandExecutorClient) SetBasicAuth(*httpoptions.BasicAuth) error {
	return tracederrors.TracedErrorNotImplemented()
}
