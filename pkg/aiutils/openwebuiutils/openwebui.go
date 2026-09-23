package openwebuiutils

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/asciich/asciichgolangpublic/pkg/httputils"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"github.com/asciich/asciichgolangpublic/pkg/urlsutils"
	"github.com/asciich/asciichgolangpublic/pkg/versionutils"
)

type OpenWebUI struct {
	Url string
}

func NewOpenWebUI(url string) (*OpenWebUI, error) {
	got := &OpenWebUI{}

	err := got.SetUrl(url)
	if err != nil {
		return nil, err
	}

	return got, err
}

func (o *OpenWebUI) SetUrl(url string) error {
	err := urlsutils.CheckIsUrl(url)
	if err != nil {
		return err
	}

	o.Url = url

	return nil
}

func (o *OpenWebUI) GetUrl() (string, error) {
	if o.Url == "" {
		return "", tracederrors.TracedError("Url not set")
	}

	return o.Url, nil
}

func (o *OpenWebUI) GetVersion(ctx context.Context) (string, error) {
	url, err := o.GetUrl()
	if err != nil {
		return "", err
	}

	// OpenWebUI typically exposes version info at /api/version
	body, err := httputils.SendRequestAndGetBodyAsBytes(
		ctx,
		&httpoptions.RequestOptions{
			Url: fmt.Sprintf("%s/api/version", url),
		},
	)
	if err != nil {
		return "", err
	}

	type Body struct {
		Version string `json:"version"`
	}

	parsed := &Body{}
	err = json.Unmarshal(body, parsed)
	if err != nil {
		return "", tracederrors.TracedErrorf("Failed to parse version response: %w", err)
	}

	version := parsed.Version

	if version == "" {
		return "", tracederrors.TracedErrorf("version is empty string after evaluation")
	}

	err = versionutils.CheckSemanticVersionString(version)
	if err != nil {
		return "", err
	}

	logging.LogInfoByCtxf(ctx, "OpenWebUI %s is running on version %s.", url, version)

	return version, nil
}

func (o *OpenWebUI) GetHealthStatus(ctx context.Context) (bool, error) {
	url, err := o.GetUrl()
	if err != nil {
		return false, err
	}

	// OpenWebUI typically has a health endpoint at /api/health
	_, err = httputils.SendRequestAndGetBodyAsBytes(
		ctx,
		&httpoptions.RequestOptions{
			Url: fmt.Sprintf("%s/api/health", url),
		},
	)
	if err != nil {
		return false, err
	}

	return true, nil
}
