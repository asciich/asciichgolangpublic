package atlassianconfluenceutils

import (
	"context"
	"errors"
	"os"
	"regexp"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/fileformats/htmlutils"
	"github.com/asciich/asciichgolangpublic/pkg/fileformats/jsonutils"
	"github.com/asciich/asciichgolangpublic/pkg/httputils"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httputilsparameteroptions"
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

// Get the content of a wiki page
func GetPageContent(ctx context.Context, url string, options *GetContentOptions) (string, error) {
	if url == "" {
		return "", tracederrors.TracedErrorEmptyString("url")
	}

	if options == nil {
		options = new(GetContentOptions)
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

	response, err := httputils.SendRequestAndGetBodyAsString(ctx, &httputilsparameteroptions.RequestOptions{
		Url:    requestUrl,
		Method: "GET",
		Header: header,
	})
	if err != nil {
		return "", err
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

	response, err := httputils.SendRequest(ctx, &httputilsparameteroptions.RequestOptions{
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

var regexPageSlashPageId = regexp.MustCompile(`\/pages\/\d+\/`)

func GetPageIdFromUrl(ctx context.Context, url string) (string, error) {
	if url == "" {
		return "", tracederrors.TracedErrorEmptyString("url")
	}

	var pageId string

	got := regexPageSlashPageId.FindString(url)
	if got != "" {
		got = strings.TrimPrefix(got, "/pages/")
		pageId = strings.TrimSuffix(got, "/")
	}

	logging.LogInfoByCtxf(ctx, "Page ID from '%s' is '%s'.", url, pageId)

	return pageId, nil
}
