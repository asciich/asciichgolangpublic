package openhandsutils

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containerinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/httputils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// Path where the CA bundle is made available inside the container.
// This is the default location on Debian/Ubuntu based images.
const containerCaBundlePath = "/etc/ssl/certs/ca-certificates.crt"

// Well known locations of the system wide CA bundle on the host.
var hostCaBundleCandidates = []string{
	"/etc/ssl/certs/ca-certificates.crt", // Debian, Ubuntu, Alpine, Arch (symlink)
	"/etc/pki/tls/certs/ca-bundle.crt",   // RHEL, CentOS, Fedora, Rocky, Alma
	"/etc/ssl/ca-bundle.pem",             // openSUSE
	"/etc/ssl/cert.pem",                  // Alpine, macOS-like layouts
}

// Returns the path of the trusted CA bundle of the host.
func getHostCaBundlePath() (string, error) {
	for _, candidate := range hostCaBundleCandidates {
		info, err := os.Stat(candidate) // follows symlinks
		if err == nil && !info.IsDir() {
			return candidate, nil
		}
	}

	return "", tracederrors.TracedErrorf(
		"Unable to find the trusted CA bundle of the host. Checked: %v",
		hostCaBundleCandidates,
	)
}

// Returns true if the given bind address (e.g. "127.0.0.1:8000") is loopback only.
func isLoopbackBindAddress(bindAddress string) bool {
	return strings.HasPrefix(bindAddress, "127.0.0.1:") ||
		strings.HasPrefix(bindAddress, "localhost:") ||
		strings.HasPrefix(bindAddress, "[::1]:")
}

// Returns the path of the file where the session API key for the given
// container is persisted on the host.
func getSessionApiKeyFilePath(containerName string) (string, error) {
	if containerName == "" {
		return "", tracederrors.TracedError("containerName is empty")
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", tracederrors.TracedErrorf("Unable to get user config dir to store the session API key: %w", err)
	}

	return filepath.Join(configDir, "asciichgolangpublic", "openhands", containerName+".apikey"), nil
}

// Returns the persisted session API key of the container.
// If none exists yet, a new random one is generated and persisted (mode 0600).
//
// The key must be persisted since an already running container (SkipIfAlreadyRunning)
// keeps the key it was started with. A new random key on every run would not match anymore.
func getOrCreateSessionApiKey(ctx context.Context, containerName string) (string, error) {
	keyFilePath, err := getSessionApiKeyFilePath(containerName)
	if err != nil {
		return "", err
	}

	content, err := os.ReadFile(keyFilePath)
	if err == nil {
		key := strings.TrimSpace(string(content))
		if key != "" {
			logging.LogInfoByCtxf(ctx, "Reuse persisted OpenHands session API key from '%s'.", keyFilePath)
			return key, nil
		}
	} else if !os.IsNotExist(err) {
		return "", tracederrors.TracedErrorf("Unable to read session API key file '%s': %w", keyFilePath, err)
	}

	key, err := GenerateSessionApiKey()
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(filepath.Dir(keyFilePath), 0o700)
	if err != nil {
		return "", tracederrors.TracedErrorf("Unable to create directory for session API key file '%s': %w", keyFilePath, err)
	}

	err = os.WriteFile(keyFilePath, []byte(key+"\n"), 0o600)
	if err != nil {
		return "", tracederrors.TracedErrorf("Unable to write session API key file '%s': %w", keyFilePath, err)
	}

	logging.LogInfoByCtxf(ctx, "Generated new OpenHands session API key and stored it in '%s'.", keyFilePath)

	return key, nil
}

func StartAsDockerContainer(ctx context.Context, options *StartContainerOptions) (containerinterfaces.Container, error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	containerName, err := options.GetContainerName()
	if err != nil {
		return nil, err
	}

	port, err := options.GetPort()
	if err != nil {
		return nil, err
	}

	bindAddress, err := options.GetBindAddress()
	if err != nil {
		return nil, err
	}

	// The session key is injected into the served page (see AllowLanSessionKeyEnvVar below).
	// This is only acceptable when the port is not reachable from the network.
	if !isLoopbackBindAddress(bindAddress) {
		return nil, tracederrors.TracedErrorf(
			"Refusing to start: bind address '%s' is not loopback but the session key is injected into the served page.",
			bindAddress,
		)
	}

	logging.LogInfoByCtxf(ctx, "Start OpenHands container '%s' bound to port %s started.", containerName, bindAddress)

	workspacePath, err := options.GetWorkspacePath()
	if err != nil {
		return nil, err
	}

	hostCaBundle, err := getHostCaBundlePath()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Mount trusted CA bundle '%s' into the container as '%s'.", hostCaBundle, containerCaBundlePath)

	sessionApiKey, err := getOrCreateSessionApiKey(ctx, containerName)
	if err != nil {
		return nil, err
	}

	container, err := dockerutils.RunContainer(ctx,
		&dockeroptions.DockerRunContainerOptions{
			Name:                 containerName,
			ImageName:            "ghcr.io/openhands/agent-canvas:sha-d92e4a0",
			KeepStoppedContainer: false,
			Ports:                []string{bindAddress + ":8000"},
			AdditionalEnvVars: map[string]string{
				"RUNTIME":       "local",
				"DEFAULT_AGENT": "CodeActAgent",

				// Fixed API key, so the Go client and the served page use the same key.
				LocalBackendApiKeyEnvVar: sessionApiKey,
				// Alternative name mentioned for the OpenHands backend.
				// Not verified, harmless if unused. Remove once confirmed.
				"OH_SESSION_API_KEYS_0": sessionApiKey,

				// Let the server inject the key into the served page again, so the
				// browser connects without the "Add a backend" dialog.
				// Only safe because the port is bound to loopback (checked above).
				AllowLanSessionKeyEnvVar: "true",

				// Make common runtimes use the mounted CA bundle.
				"SSL_CERT_FILE":       containerCaBundlePath, // Go, OpenSSL, Python ssl
				"REQUESTS_CA_BUNDLE":  containerCaBundlePath, // Python requests
				"CURL_CA_BUNDLE":      containerCaBundlePath, // curl
				"NODE_EXTRA_CA_CERTS": containerCaBundlePath, // Node.js
			},
			Mounts: []string{
				workspacePath + ":/workspace",
				hostCaBundle + ":" + containerCaBundlePath + ":ro",
			},
			SkipIfAlreadyRunning: true,
		},
	)
	if err != nil {
		return nil, err
	}

	url := "http://127.0.0.1:" + strconv.Itoa(port)

	err = httputils.WaitUntilStatusCodeOK(ctx, url, time.Minute*2)
	if err != nil {
		return nil, err
	}

	openhands, err := NewOpenHands(url)
	if err != nil {
		return nil, err
	}

	err = openhands.SetSessionApiKey(sessionApiKey)
	if err != nil {
		return nil, err
	}

	err = openhands.DeclineAnalytics(ctx)
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(
		ctx,
		"OpenHands API key (only needed if the 'Add a backend' dialog appears at %s/canvas/) is stored in the key file. Use: cat ~/.config/asciichgolangpublic/openhands/%s.apikey",
		url,
		containerName,
	)

	logging.LogInfoByCtxf(ctx, "Going to add the default workspace.")
	err = openhands.CreateWorkspace(ctx, "workspace", "/workspace")
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Start OpenHands container '%s' bound to port %s finished.", containerName, bindAddress)

	return container, nil
}
