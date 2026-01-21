package atlassianconfluenceutils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/datatypes/slicesutils"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/documentutils/documentconverter"
	"github.com/asciich/asciichgolangpublic/pkg/fileformats/htmlutils"
	"github.com/asciich/asciichgolangpublic/pkg/fileformats/jsonutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/httputils"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"github.com/asciich/asciichgolangpublic/pkg/urlsutils"
)

const TokenEnvVarName = "ATLASSIAN_CONFLUENCE_TOKEN"
const CookieEnvVarName = "ATLASSIAN_CONFLUENCE_COOKIE"

var ErrTokenNotSet = errors.New("atlassian Confluence wiki token not set in the env var '" + TokenEnvVarName + "'")
var ErrCookieNotSet = errors.New("atlassian Confluence wiki cookie not set in the env var '" + CookieEnvVarName + "'")

func GetToken(ctx context.Context) (string, error) {
	token := os.Getenv(TokenEnvVarName)
	if token == "" {
		return "", tracederrors.TracedError(ErrTokenNotSet)
	} else {
		logging.LogInfoByCtxf(ctx, "Loaded atlassian token from env var '%s'.", TokenEnvVarName)
	}

	return token, nil
}

func GetCookieFromEnvVar(ctx context.Context) (string, error) {
	cookie := os.Getenv(CookieEnvVarName)
	cookie = strings.TrimSpace(cookie)
	if cookie == "" {
		return "", tracederrors.TracedError(ErrCookieNotSet)
	} else {
		logging.LogInfoByCtxf(ctx, "Loaded atlassian cookie from env var '%s'.", CookieEnvVarName)
	}

	return cookie, nil
}

// Get a map[string]string of Headers to use in a Request towards the conflunece wiki
func getRequestHeader(ctx context.Context) (map[string]string, error) {
	header := map[string]string{
		"Content-Type": "application/json",
	}

	var authenticationLoaded bool

	token, err := GetToken(ctx)
	if err == nil {
		header["Authorization"] = "Bearer " + token
		logging.LogInfoByCtxf(ctx, "Use Atlassian confluence token in request header.")
		authenticationLoaded = true
	} else {
		if !errors.Is(err, ErrTokenNotSet) {
			return nil, err
		}
	}

	cookie, err := GetCookieFromEnvVar(ctx)
	if err == nil {
		header["Cookie"] = cookie
		logging.LogInfoByCtxf(ctx, "Use Atlassian confluence cookie in request header.")
		authenticationLoaded = true
	} else {
		if !errors.Is(err, ErrCookieNotSet) {
			return nil, err
		}
	}

	if !authenticationLoaded {
		return nil, tracederrors.TracedErrorf("Unable to get request header. No authentication method found. Please set at least one of these environment variables: '%s' or '%s'.", TokenEnvVarName, CookieEnvVarName)
	}

	return header, nil
}

// Returns a slice of the page Id's of all subpages:
func GetChildPageIds(ctx context.Context, url string, options *GetChildPageOptions) ([]string, error) {
	if url == "" {
		return nil, tracederrors.TracedErrorEmptyString("url")
	}

	if options == nil {
		options = new(GetChildPageOptions)
	}

	pageId, err := GetPageIdFromUrl(ctx, url)
	if err != nil {
		return nil, err
	}

	header, err := getRequestHeader(ctx)
	if err != nil {
		return nil, err
	}

	baseUrl, err := urlsutils.GetBaseUrl(url)
	if err != nil {
		return nil, err
	}

	requestUrl := baseUrl + "/rest/api/content/" + pageId + "/child/page"
	logging.LogInfoByCtxf(ctx, "Use atlassian confluence api url %s to retrieve the pages child pages.", requestUrl)

	response, err := httputils.SendRequestAndGetBodyAsString(ctx, &httpoptions.RequestOptions{
		Url:    requestUrl,
		Method: "GET",
		Header: header,
	})
	if err != nil {
		return nil, err
	}

	allIds, err := jsonutils.RunJqAgainstJsonStringAsString(response, ".results[].id")
	if err != nil {
		return nil, err
	}

	ids := stringsutils.SplitLines(allIds, true)

	if options.Recursive {
		logging.LogInfoByCtxf(ctx, "Going to recursively collect child pages.")

		allIds := slicesutils.GetDeepCopyOfStringsSlice(ids)
		for _, id := range ids {
			subIds, err := GetChildPageIds(ctx, baseUrl+"/rest/api/content/"+id+"/child/page", options)
			if err != nil {
				return nil, err
			}

			allIds = append(allIds, subIds...)
		}

		ids = allIds
	}

	logging.LogInfoByCtxf(ctx, "Found '%d' subpages for the wiki page '%s'", len(ids), pageId)

	return ids, nil
}

// Get the raw response of the API for the wiki page body.storage:
func GetBodyStorageRawApiResponse(ctx context.Context, url string) (string, error) {
	if url == "" {
		return "", tracederrors.TracedErrorEmptyString("url")
	}

	pageId, err := GetPageIdFromUrl(ctx, url)
	if err != nil {
		return "", err
	}

	header, err := getRequestHeader(ctx)
	if err != nil {
		return "", err
	}

	baseUrl, err := urlsutils.GetBaseUrl(url)
	if err != nil {
		return "", err
	}

	requestUrl := baseUrl + "/rest/api/content/" + pageId + "?expand=body.storage"
	logging.LogInfoByCtxf(ctx, "Use atlassian confluence api url %s to retrieve the page content.", requestUrl)

	response, err := httputils.SendRequestAndGetBodyAsString(ctx, &httpoptions.RequestOptions{
		Url:    requestUrl,
		Method: "GET",
		Header: header,
	})
	if err != nil {
		return "", err
	}

	return response, nil
}

// Get the content of a wiki page
func GetPageContent(ctx context.Context, url string, options *GetContentOptions) (string, error) {
	response, err := GetBodyStorageRawApiResponse(ctx, url)
	if err != nil {
		return "", err
	}

	if options == nil {
		options = new(GetContentOptions)
	}

	body, err := jsonutils.RunJqAgainstJsonStringAsString(response, ".body.storage.value")
	if err != nil {
		return "", err
	}

	if options.PrettyPrint {
		body, err = htmlutils.PrettyFormat(body)
		if err != nil {
			return "", err
		}
	}

	return body, nil
}

func GetRequest(ctx context.Context, url string) (string, error) {
	if url == "" {
		return "", tracederrors.TracedErrorEmptyString("url")
	}

	header, err := getRequestHeader(ctx)
	if err != nil {
		return "", err
	}

	response, err := httputils.SendRequest(ctx, &httpoptions.RequestOptions{
		Url:    url,
		Method: "GET",
		Header: header,
	})
	if err != nil {
		return "", err
	}

	body, err := response.GetBodyAsString()
	if err != nil {
		return "", err
	}

	return body, nil
}

var regexPageIdBetweenSlashes = regexp.MustCompile(`\/\d+\/`)
var regexPageId = regexp.MustCompile(`^\d+$`)

func GetPageIdFromUrl(ctx context.Context, url string) (string, error) {
	if url == "" {
		return "", tracederrors.TracedErrorEmptyString("url")
	}

	pageId := regexPageIdBetweenSlashes.FindString(url)
	pageId = strings.TrimPrefix(pageId, "/")
	pageId = strings.TrimSuffix(pageId, "/")

	if pageId == "" {
		pageId = regexPageId.FindString(url)
	}

	if pageId == "" {
		return "", tracederrors.TracedErrorf("Unable to get pageId from given url='%s'. pageId is emtpy string after evaluation.", url)
	}

	logging.LogInfoByCtxf(ctx, "Page ID from '%s' is '%s'.", url, pageId)

	return pageId, nil
}

func downloadSinglePageContent(ctx context.Context, url string, outputDir string, options *DownloadPageContentOptions) error {
	if options == nil {
		options = &DownloadPageContentOptions{}
	}

	if url == "" {
		return tracederrors.TracedErrorEmptyString("url")
	}

	if outputDir == "" {
		return tracederrors.TracedErrorEmptyString("outputDir")
	}

	baseUrl, err := urlsutils.GetBaseUrl(url)
	if err != nil {
		return err
	}

	pageId, err := GetPageIdFromUrl(ctx, url)
	if err != nil {
		return err
	}

	rawResponse, err := GetBodyStorageRawApiResponse(ctx, url)
	if err != nil {
		return err
	}

	webuiPath, err := jsonutils.RunJqAgainstJsonStringAsString(rawResponse, "._links.webui")
	if err != nil {
		return err
	}
	if webuiPath == "" {
		return tracederrors.TracedError("webuiPath is empty string after evaluation.")
	}

	pageUrl, err := jsonutils.RunJqAgainstJsonStringAsString(rawResponse, "._links.self")
	if err != nil {
		return err
	}
	if pageUrl == "" {
		return tracederrors.TracedError("pageUrl is empty string after evaluation.")
	}

	pageTitle, err := jsonutils.RunJqAgainstJsonStringAsString(rawResponse, ".title")
	if err != nil {
		return err
	}

	pageOutputDir := filepath.Join(outputDir, strings.ReplaceAll(strings.ReplaceAll(baseUrl, "/", "_"), ":", "_"), webuiPath)

	err = nativefiles.CreateDirectory(ctx, pageOutputDir, &filesoptions.CreateOptions{})
	if err != nil {
		return err
	}

	contentBaseName := fmt.Sprintf("%s_content.html", pageId)
	if options.ConvertToMdFiles {
		contentBaseName = fmt.Sprintf("%s_content.md", pageId)
	}

	content, err := jsonutils.RunJqAgainstJsonStringAsString(rawResponse, ".body.storage.value")
	if err != nil {
		return err
	}

	if options.ConvertToMdFiles {
		content, err = documentconverter.HtmlStringToMdString(content)
		if err != nil {
			return err
		}
	} else {
		content, err = htmlutils.PrettyFormat(content)
		if err != nil {
			return err
		}
	}

	err = nativefiles.WriteString(ctx, filepath.Join(pageOutputDir, contentBaseName), content)
	if err != nil {
		return err
	}

	type infoContent struct {
		WikiInstance string `json:"wiki-instance"`
		PageUrl      string `json:"page-url"`
		WebUiUrl     string `json:"webui-url"`
		PageId       string `json:"page-id"`
		BaseName     string `json:"basename"`
		Title        string `json:"title"`
	}

	info, err := json.MarshalIndent(&infoContent{
		WikiInstance: baseUrl,
		PageUrl:      pageUrl,
		PageId:       pageId,
		BaseName:     contentBaseName,
		WebUiUrl:     baseUrl + webuiPath,
		Title:        pageTitle,
	}, "", "  ")
	if err != nil {
		return tracederrors.TracedErrorf("Failed to marshal info file: %w", err)
	}

	err = nativefiles.WriteBytes(ctx, filepath.Join(pageOutputDir, "info.json"), info)
	if err != nil {
		return err
	}

	return nil
}

func DownloadPageContent(ctx context.Context, url string, outputDir string, options *DownloadPageContentOptions) error {
	if url == "" {
		return tracederrors.TracedErrorEmptyString("url")
	}

	if outputDir == "" {
		return tracederrors.TracedErrorEmptyString("outputDir")
	}

	logging.LogInfoByCtxf(ctx, "Download page content of %s started.", url)

	err := downloadSinglePageContent(ctx, url, outputDir, options)
	if err != nil {
		return err
	}

	if options.Recursive {
		logging.LogInfoByCtxf(ctx, "Going to collect and download the child pages of %s .", url)

		baseUrl, err := urlsutils.GetBaseUrl(url)
		if err != nil {
			return err
		}

		ids, err := GetChildPageIds(ctx, url, &GetChildPageOptions{Recursive: true})
		if err != nil {
			return err
		}

		for _, id := range ids {
			err = downloadSinglePageContent(ctx, baseUrl+"/rest/api/content/"+id+"/child/page", outputDir, options)
			if err != nil {
				return err
			}
		}
	}

	logging.LogInfoByCtxf(ctx, "Download page content of %s finished.", url)

	return nil
}
