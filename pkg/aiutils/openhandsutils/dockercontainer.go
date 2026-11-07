package openhandsutils

import (
	"context"
	"os"
	"strconv"
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
	"/etc/ssl/certs/ca-certificates.crt",   // Debian, Ubuntu, Alpine, Arch (symlink)
	"/etc/pki/tls/certs/ca-bundle.crt",     // RHEL, CentOS, Fedora, Rocky, Alma
	"/etc/ssl/ca-bundle.pem",               // openSUSE
	"/etc/ssl/cert.pem",                    // Alpine, macOS-like layouts
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

	container, err := dockerutils.RunContainer(ctx,
		&dockeroptions.DockerRunContainerOptions{
			Name:                 containerName,
			ImageName:            "ghcr.io/openhands/agent-canvas:latest",
			KeepStoppedContainer: false,
			Ports:                []string{bindAddress + ":8000"},
			AdditionalEnvVars: map[string]string{
				"RUNTIME":       "local",
				"DEFAULT_AGENT": "CodeActAgent",

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

	logging.LogInfoByCtxf(ctx, "Going to add the default workspace.")
	err = openhands.CreateWorkspace(ctx, "workspace", "/workspace")
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Start OpenHands container '%s' bound to port %s finished.", containerName, bindAddress)

	return container, nil
}
