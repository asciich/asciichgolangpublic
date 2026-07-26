package openhandsutils

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/httputils"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"github.com/asciich/asciichgolangpublic/pkg/urlsutils"
	"github.com/asciich/asciichgolangpublic/pkg/versionutils"
)

type Openhands struct {
	Url string
}

func NewOpenHands(url string) (*Openhands, error) {
	got := &Openhands{}

	err := got.SetUrl(url)
	if err != nil {
		return nil, err
	}

	return got, err
}

func (o *Openhands) SetUrl(url string) error {
	err := urlsutils.CheckIsUrl(url)
	if err != nil {
		return err
	}

	o.Url = url

	return nil
}

func (o *Openhands) GetUrl() (string, error) {
	if o.Url == "" {
		return "", tracederrors.TracedError("Url not set")
	}

	return o.Url, nil
}

func (o *Openhands) GetServerInfoRawResponse(ctx context.Context) ([]byte, error) {
	url, err := o.GetUrl()
	if err != nil {
		return nil, err
	}

	return httputils.SendRequestAndGetBodyAsBytes(
		ctx,
		&httpoptions.RequestOptions{
			Url: fmt.Sprintf("%s/server_info", url),
		},
	)
}

func (o *Openhands) GetVersion(ctx context.Context) (string, error) {
	url, err := o.GetUrl()
	if err != nil {
		return "", err
	}

	serverInfoBody, err := o.GetServerInfoRawResponse(ctx)
	if err != nil {
		return "", err
	}

	type Body struct {
		Version string `json:"version"`
	}

	parsed := &Body{}
	err = json.Unmarshal(serverInfoBody, parsed)
	if err != nil {
		return "", tracederrors.TracedErrorf("Failed to parse server raw response to get version: %w", err)
	}

	version := parsed.Version

	if version == "" {
		return "", tracederrors.TracedErrorf("version is empty string after evaluation")
	}

	err = versionutils.CheckSemanticVersionString(version)
	if err != nil {
		return "", err
	}

	logging.LogInfoByCtxf(ctx, "OpenHands %s is running on version %s.", url, version)

	return version, nil
}

func (o *Openhands) GetSessionApiKey(ctx context.Context) (string, error) {
	url, err := o.GetUrl()
	if err != nil {
		return "", err
	}

	body, err := httputils.SendRequestAndGetBodyAsBytes(
		ctx,
		&httpoptions.RequestOptions{
			Url: url,
		},
	)
	if err != nil {
		return "", err
	}

	// Extract the key from: window.__AGENT_CANVAS_SESSION_API_KEY__="<key>"
	marker := `window.__AGENT_CANVAS_SESSION_API_KEY__="`
	content := string(body)

	idx := strings.Index(content, marker)
	if idx == -1 {
		return "", tracederrors.TracedError("session API key not found in page")
	}

	start := idx + len(marker)
	end := strings.Index(content[start:], `"`)
	if end == -1 {
		return "", tracederrors.TracedError("could not parse session API key")
	}

	apiKey := content[start : start+end]

	if len(apiKey) <= 0 {
		return "", tracederrors.TracedError("Unable to extract valid API key")
	}

	return apiKey, nil
}

func (o *Openhands) DoApiRequest(ctx context.Context, method string, path string, payload []byte) ([]byte, error) {
	sessionKey, err := o.GetSessionApiKey(ctx)
	if err != nil {
		return nil, err
	}

	url, err := o.GetUrl()
	if err != nil {
		return nil, err
	}

	return httputils.SendRequestAndGetBodyAsBytes(
		ctx,
		&httpoptions.RequestOptions{
			Url:    url + path,
			Method: method,
			Header: map[string]string{
				"X-Session-API-Key": sessionKey,
				"Content-Type":      "application/json",
			},
			Data: payload,
		},
	)
}
