package openhandsutils

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/httputils"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"github.com/asciich/asciichgolangpublic/pkg/urlsutils"
	"github.com/asciich/asciichgolangpublic/pkg/versionutils"
)

// Environment variable used by newer OpenHands / Agent Canvas versions to
// configure the API key of the backend.
const LocalBackendApiKeyEnvVar = "LOCAL_BACKEND_API_KEY"

// Environment variable which makes newer versions inject the session key into
// the served index page again. Only use it together with a port which is bound
// to loopback (127.0.0.1) only.
const AllowLanSessionKeyEnvVar = "AGENT_CANVAS_ALLOW_LAN_SESSION_KEY"

// Matches: window.__AGENT_CANVAS_SESSION_API_KEY__="<key>"
// Also tolerates whitespace around "=" and single quotes.
var sessionApiKeyRegex = regexp.MustCompile(
	`window\.__AGENT_CANVAS_SESSION_API_KEY__\s*=\s*("(?:[^"\\]|\\.)*"|'[^']*')`,
)

type Openhands struct {
	Url string

	// Optional: explicitly configured session API key, for example
	//   - the value passed to the container as LOCAL_BACKEND_API_KEY
	//   - the content of $STATE_DIR/api-key.txt inside the container
	// If set, the key is NOT scraped from the web page.
	SessionApiKey string
}

func NewOpenHands(url string) (*Openhands, error) {
	got := &Openhands{}

	err := got.SetUrl(url)
	if err != nil {
		return nil, err
	}

	return got, err
}

// Generates a random 256 bit session API key, hex encoded.
// Pass it to the container as LOCAL_BACKEND_API_KEY and set it with SetSessionApiKey.
func GenerateSessionApiKey() (string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return "", tracederrors.TracedErrorf("Failed to generate random session API key: %w", err)
	}

	return hex.EncodeToString(b), nil
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

func (o *Openhands) SetSessionApiKey(key string) error {
	key = strings.TrimSpace(key)
	if key == "" {
		return tracederrors.TracedError("session API key is empty")
	}

	o.SessionApiKey = key

	return nil
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

// Returns the session API key.
//
// Order:
//  1. Explicitly configured key (o.SessionApiKey)
//  2. Key injected into the index page by the server (legacy / opt-in behaviour)
//
// Newer versions do not inject the key into the page by default anymore.
// Prefer setting the key explicitly via SetSessionApiKey.
func (o *Openhands) GetSessionApiKey(ctx context.Context) (string, error) {
	if o.SessionApiKey != "" {
		return o.SessionApiKey, nil
	}

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

	matches := sessionApiKeyRegex.FindSubmatch(body)
	if matches == nil {
		return "", tracederrors.TracedErrorf(
			"session API key not found in page. Newer OpenHands versions do not inject the key by default. "+
				"Either set the key explicitly (start the container with '-e %s=<key>' and call SetSessionApiKey), "+
				"or start the container with '-e %s=true' and a loopback-only port binding.",
			LocalBackendApiKeyEnvVar,
			AllowLanSessionKeyEnvVar,
		)
	}

	literal := string(matches[1])

	var apiKey string
	if strings.HasPrefix(literal, `"`) {
		// JSON compatible string literal -> unescape properly
		err = json.Unmarshal([]byte(literal), &apiKey)
		if err != nil {
			return "", tracederrors.TracedErrorf("could not parse session API key: %w", err)
		}
	} else {
		apiKey = strings.Trim(literal, `'`)
	}

	if apiKey == "" {
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
				// Sent in both styles because I could not confirm which one
				// the new backend expects. You can drop one once verified.
				"X-Session-API-Key": sessionKey,
				"Authorization":     "Bearer " + sessionKey,
				"Content-Type":      "application/json",
			},
			Data: payload,
		},
	)
}

// Declines the anonymous usage data collection so the
// "Help improve OpenHands" dialog does not show up.
//
// The setting lives in misc_settings.app_preferences.user_consents_to_analytics.
// An existing value (e.g. chosen by the user in the UI) is never overwritten.
//
// PATCH /api/settings expects "misc_settings_diff". It is not known whether the diff is
// merged deeply or only on the first level. Therefore the complete app_preferences
// object is sent, so no other preference (e.g. enabled_skills) can get lost.
func (o *Openhands) DeclineAnalytics(ctx context.Context) error {
	const field = "user_consents_to_analytics"

	body, err := o.DoApiRequest(ctx, "GET", "/api/settings", nil)
	if err != nil {
		return err
	}

	settings := map[string]any{}
	err = json.Unmarshal(body, &settings)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to parse OpenHands settings: %w", err)
	}

	miscSettings, _ := settings["misc_settings"].(map[string]any)

	appPreferences, _ := miscSettings["app_preferences"].(map[string]any)
	if appPreferences == nil {
		appPreferences = map[string]any{}
	}

	if existing, ok := appPreferences[field]; ok {
		logging.LogInfoByCtxf(ctx, "OpenHands analytics consent is already set to '%v'. Leave it unchanged.", existing)
		return nil
	}

	appPreferences[field] = false

	payload, err := json.Marshal(map[string]any{
		"misc_settings_diff": map[string]any{
			"app_preferences": appPreferences,
		},
	})
	if err != nil {
		return tracederrors.TracedErrorf("Failed to marshal analytics settings: %w", err)
	}

	_, err = o.DoApiRequest(ctx, "PATCH", "/api/settings", payload)
	if err != nil {
		return err
	}

	// Verify, so a silently ignored PATCH does not go unnoticed.
	body, err = o.DoApiRequest(ctx, "GET", "/api/settings", nil)
	if err != nil {
		return err
	}

	type verifyBody struct {
		MiscSettings struct {
			AppPreferences map[string]any `json:"app_preferences"`
		} `json:"misc_settings"`
	}

	verify := &verifyBody{}
	err = json.Unmarshal(body, verify)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to parse OpenHands settings for verification: %w", err)
	}

	value, ok := verify.MiscSettings.AppPreferences[field]
	if !ok || value != false {
		return tracederrors.TracedErrorf(
			"Declining analytics had no effect. '%s' is '%v' after the PATCH.", field, value,
		)
	}

	// Make sure the other preferences survived the PATCH.
	if _, ok := verify.MiscSettings.AppPreferences["enabled_skills"]; !ok {
		logging.LogWarnByCtxf(ctx, "'enabled_skills' is missing in app_preferences after the PATCH. The diff may have replaced the object.")
	}

	logging.LogInfoByCtxf(ctx, "Declined anonymous usage data collection of OpenHands.")

	return nil
}
