package pfsenseutils

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

const ENV_VAR_NAME_PFSENSE_PASSWORD = "PFSENSE_PASSWORD"

type Router struct {
	Url      string
	UserName string
	Password string

	client *http.Client
}

func (r *Router) GetUrl() (string, error) {
	if r.Url == "" {
		return "", tracederrors.TracedError("Url not set")
	}

	return r.Url, nil
}

func (r *Router) GetUserName() (string, error) {
	username := r.UserName

	if username == "" {
		username = os.Getenv("PFSENSE_USER")
	}

	if username == "" {
		return "", tracederrors.TracedError("username not set")
	}

	return username, nil
}

func (r *Router) GetPassword() (string, error) {
	password := r.Password

	if password == "" {
		password = os.Getenv(ENV_VAR_NAME_PFSENSE_PASSWORD)
	}

	if password == "" {
		return "", tracederrors.TracedErrorf("password not set. Set the '%s' env var.", ENV_VAR_NAME_PFSENSE_PASSWORD)
	}

	return password, nil
}

func (r *Router) IsLoggedIn(ctx context.Context) (bool, error) {
	url, err := r.GetUrl()
	if err != nil {
		return false, err
	}

	logging.LogInfoByCtxf(ctx, "Is logged in to pfSense router '%s' started.", url)

	var isLoggedIn bool
	if r.client != nil {
		resp, err := r.client.Get(url)
		if err != nil {
			r.client = nil
		} else {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			dashboardHTML := string(body)
			if strings.Contains(dashboardHTML, "Sign In") && !strings.Contains(dashboardHTML, "Dashboard") {
				r.client = nil
			}
			isLoggedIn = true
		}
	}

	if isLoggedIn {
		logging.LogInfoByCtxf(ctx, "Logged in to pfSense router '%s'.", url)
	} else {
		logging.LogInfoByCtxf(ctx, "Not logged in to pfSense router '%s'.", url)
	}

	logging.LogInfoByCtxf(ctx, "Is logged in to pfSense router '%s' finished.", url)

	return isLoggedIn, nil
}

func (r *Router) Login(ctx context.Context) error {
	pfsenseUrl, err := r.GetUrl()
	if err != nil {
		return err
	}

	username, err := r.GetUserName()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Login to pfsense router '%s' as user '%s' started.", pfsenseUrl, username)

	password, err := r.GetPassword()
	if err != nil {
		return err
	}

	isLoggedIn, err := r.IsLoggedIn(ctx)
	if err != nil {
		return err
	}

	if isLoggedIn {
		logging.LogInfoByCtxf(ctx, "Already logged in to pfSense router '%s'. Skip relogin.", pfsenseUrl)
	} else {
		jar, _ := cookiejar.New(nil)
		r.client = &http.Client{
			Jar: jar,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			// Don't follow redirects automatically so we can inspect responses
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}

		logging.LogInfoByCtxf(ctx, "Fetch pfSense login page %s .", pfsenseUrl)
		resp, err := r.client.Get(pfsenseUrl)
		if err != nil {
			return tracederrors.TracedErrorf("Failed to fetch login page %s : %w", pfsenseUrl, err)
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		csrfToken := extractCSRFToken(string(body))
		if csrfToken == "" {
			return tracederrors.TracedError("Failed to extract CSRF token from login page")
		}

		logging.LogInfoByCtxf(ctx, "Perform pfSense login.")
		formData := url.Values{
			"__csrf_magic": {csrfToken},
			"usernamefld":  {username},
			"passwordfld":  {password},
			"login":        {"Sign In"},
		}

		resp, err = r.client.PostForm(pfsenseUrl, formData)
		if err != nil {
			return tracederrors.TracedErrorf("Failed to submit login form: %w", err)
		}
		resp.Body.Close()

		isLoggedIn, err := r.IsLoggedIn(ctx)
		if err != nil {
			return err
		}

		if !isLoggedIn {
			return tracederrors.TracedErrorf("Login failed. Check for login returnes false after login.")
		}

		logging.LogChangedByCtxf(ctx, "Successfully logged in to pfSense '%s'.", pfsenseUrl)
	}

	logging.LogInfoByCtxf(ctx, "Login to pfSense router '%s' as user '%s' finished.", pfsenseUrl, username)

	return err
}

// extractCSRFToken finds the CSRF magic token in the HTML response
func extractCSRFToken(html string) string {
	// pfSense uses a hidden input field named "__csrf_magic"
	re := regexp.MustCompile(`name="__csrf_magic"\s+value="([^"]+)"`)
	matches := re.FindStringSubmatch(html)
	if len(matches) > 1 {
		return matches[1]
	}

	// Alternative pattern
	re = regexp.MustCompile(`var\s+csrfMagicToken\s*=\s*"([^"]+)"`)
	matches = re.FindStringSubmatch(html)
	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}
