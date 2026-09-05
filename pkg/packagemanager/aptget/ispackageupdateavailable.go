package aptget

import (
	"context"
	"regexp"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/packagemanageroptions"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func IsPackagesUpdateAvailable(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, packageNames []string, options *packagemanageroptions.UpdateDatabaseOptions) (bool, error) {
	if len(packageNames) == 0 {
		return false, tracederrors.TracedError("packageNames has no entries.")
	}

	logging.LogInfoByCtxf(ctx, "Is apt-get packages update available for '%s' started.", packageNames)
	var updateAvailable bool
	var err error

	for _, name := range packageNames {
		updateAvailable, err = IsPackageUpdateAvailable(ctx, commandExecutor, name, options)
		if err != nil {
			return false, err
		}

		if updateAvailable {
			updateAvailable = true
			break
		}
	}

	logging.LogInfoByCtxf(ctx, "Is apt-get packages update available for '%s' finished.", packageNames)

	return updateAvailable, nil
}

// parseInstalledAndCandidate extracts the "Installed:" and "Candidate:" versions
// from the output of "apt-cache policy <package>".
func parseInstalledAndCandidate(policyOutput string) (installed string, candidate string) {
	installedRegex := regexp.MustCompile(`(?m)^\s*Installed:\s*(.+)\s*$`)
	candidateRegex := regexp.MustCompile(`(?m)^\s*Candidate:\s*(.+)\s*$`)

	if match := installedRegex.FindStringSubmatch(policyOutput); match != nil {
		installed = strings.TrimSpace(match[1])
	}
	if match := candidateRegex.FindStringSubmatch(policyOutput); match != nil {
		candidate = strings.TrimSpace(match[1])
	}

	return installed, candidate
}
func IsPackageUpdateAvailable(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, packageName string, options *packagemanageroptions.UpdateDatabaseOptions) (bool, error) {
	if commandExecutor == nil {
		return false, tracederrors.TracedErrorNil("commandExecutor")
	}

	if packageName == "" {
		return false, tracederrors.TracedErrorEmptyString("packageName")
	}

	logging.LogInfoByCtxf(ctx, "Is apt-get update available for package '%s' started.", packageName)

	if options == nil {
		options = new(packagemanageroptions.UpdateDatabaseOptions)
	}

	// apt has no direct equivalent of "pacman -Qu". We compare the installed
	// version against the candidate version reported by "apt-cache policy".
	queryPackagePolicy := func() (string, error) {
		output, err := commandExecutor.RunCommand(
			ctx,
			&parameteroptions.RunCommandOptions{
				Command:           []string{"apt-cache", "policy", packageName},
				AllowAllExitCodes: true,
			},
		)
		if err != nil {
			return "", err
		}

		exitCode, err := output.GetReturnCode()
		if err != nil {
			return "", err
		}
		if exitCode != 0 {
			stderr, err := output.GetStderrAsString()
			if err != nil {
				return "", err
			}
			return "", tracederrors.TracedErrorf("Failed to evaluate if an update of the apt-get package '%s' is available: %s", packageName, stderr)
		}

		return output.GetStdoutAsString()
	}

	stdout, err := queryPackagePolicy()
	if err != nil {
		return false, err
	}

	installed, candidate := parseInstalledAndCandidate(stdout)

	// The package is not installed at all, so there is nothing to update:
	if installed == "" || installed == "(none)" {
		logging.LogInfoByCtxf(ctx, "Apt-get package '%s' is not installed, no update available.", packageName)
		return false, nil
	}

	// On a fresh system the apt lists were never downloaded, so the candidate
	// is derived solely from the local dpkg status and therefore equals the
	// installed version. In that ambiguous case we refresh the database once
	// (mirroring pacman's "-Sy" fallback) and re-query to get the real
	// candidate from the repositories.
	if candidate == "" || candidate == "(none)" || candidate == installed {
		logging.LogInfoByCtxf(ctx, "Apt-get database possibly not refreshed yet for package '%s'. Going to download apt-get database and re-check.", packageName)

		err := UpdateDatabase(ctx, commandExecutor, options)
		if err != nil {
			return false, err
		}

		stdout, err = queryPackagePolicy()
		if err != nil {
			return false, err
		}

		installed, candidate = parseInstalledAndCandidate(stdout)
	}

	if candidate == "" || candidate == "(none)" {
		logging.LogInfoByCtxf(ctx, "No candidate for apt-get package '%s' available, no update available.", packageName)
		return false, nil
	}

	if installed != candidate {
		logging.LogInfoByCtxf(ctx, "Update for the package '%s' found using apt-get (installed: '%s', candidate: '%s').", packageName, installed, candidate)
		return true, nil
	}

	logging.LogInfoByCtxf(ctx, "No update for apt-get package '%s' available.", packageName)
	return false, nil
}

func (a *AptGet) IsPackageUpdateAvailable(ctx context.Context, packageName string, options *packagemanageroptions.UpdateDatabaseOptions) (bool, error) {
	commandExecutor, err := a.GetCommandExecutor()
	if err != nil {
		return false, err
	}

	return IsPackageUpdateAvailable(ctx, commandExecutor, packageName, options)
}
