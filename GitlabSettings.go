package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type GitlabSettings struct {
	gitlab *GitlabInstance
}

func NewGitlabSettings() (gitlabSettings *GitlabSettings) {
	return new(GitlabSettings)
}

func (g *GitlabSettings) MustDisableAutoDevops(verbose bool) {
	err := g.DisableAutoDevops(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabSettings) MustDisableSignup(verbose bool) {
	err := g.DisableSignup(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabSettings) MustGetCurrentSettingsNative() (nativeSettings *gitlab.Settings) {
	nativeSettings, err := g.GetCurrentSettingsNative()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nativeSettings
}

func (g *GitlabSettings) MustGetFqdn() (fqdn string) {
	fqdn, err := g.GetFqdn()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return fqdn
}

func (g *GitlabSettings) MustGetGitlab() (gitlab *GitlabInstance) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlab
}

func (g *GitlabSettings) MustGetNativeClient() (nativeClient *gitlab.Client) {
	nativeClient, err := g.GetNativeClient()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nativeClient
}

func (g *GitlabSettings) MustGetNativeSettingsService() (nativeSettingsService *gitlab.SettingsService) {
	nativeSettingsService, err := g.GetNativeSettingsService()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nativeSettingsService
}

func (g *GitlabSettings) MustIsAutoDevopsEnabled() (isAutoDevopsEnabled bool) {
	isAutoDevopsEnabled, err := g.IsAutoDevopsEnabled()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isAutoDevopsEnabled
}

func (g *GitlabSettings) MustIsSignupEnabled() (isSignupEnabled bool) {
	isSignupEnabled, err := g.IsSignupEnabled()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isSignupEnabled
}

func (g *GitlabSettings) MustSetGitlab(gitlab *GitlabInstance) {
	err := g.SetGitlab(gitlab)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (s *GitlabSettings) DisableAutoDevops(verbose bool) (err error) {
	isAutoDevopsEnabled, err := s.IsAutoDevopsEnabled()
	if err != nil {
		return err
	}

	fqdn, err := s.GetFqdn()
	if err != nil {
		return err
	}

	if isAutoDevopsEnabled {
		if verbose {
			logging.LogInfof("Disable AutoDevops on gitlab '%s' started.", fqdn)
		}

		nativeSettings, err := s.GetNativeSettingsService()
		if err != nil {
			return err
		}

		autoDevopsEnabled := false
		_, _, err = nativeSettings.UpdateSettings(&gitlab.UpdateSettingsOptions{
			AutoDevOpsEnabled: &autoDevopsEnabled,
		})
		if err != nil {
			return err
		}
		if verbose {
			logging.LogChangedf("Dislabed AutoDevops on gitlab '%s'.", fqdn)
		}

		if verbose {
			logging.LogInfof("Disable AutoDevops on gitlab '%s' finished.", fqdn)
		}
	} else {
		if verbose {
			logging.LogInfof("Autodevops on gitlab '%s' already disabled.", fqdn)
		}
	}

	return nil
}

func (s *GitlabSettings) DisableSignup(verbose bool) (err error) {
	isSignupEnabled, err := s.IsSignupEnabled()
	if err != nil {
		return err
	}

	fqdn, err := s.GetFqdn()
	if err != nil {
		return err
	}

	if isSignupEnabled {
		if verbose {
			logging.LogInfof("Disable signup to gitlab '%s' started.", fqdn)
		}

		nativeSettings, err := s.GetNativeSettingsService()
		if err != nil {
			return err
		}

		signupEnabled := false
		_, _, err = nativeSettings.UpdateSettings(&gitlab.UpdateSettingsOptions{
			SignupEnabled: &signupEnabled,
		})
		if err != nil {
			return err
		}
		if verbose {
			logging.LogChangedf("Dislabed signup to gitlab '%s'.", fqdn)
		}

		if verbose {
			logging.LogInfof("Disable signup to gitlab '%s' finished.", fqdn)
		}
	} else {
		if verbose {
			logging.LogInfof("Signup to gitlab '%s' already disabled.", fqdn)
		}
	}
	return nil
}

func (s *GitlabSettings) GetCurrentSettingsNative() (nativeSettings *gitlab.Settings, err error) {
	settings, err := s.GetNativeSettingsService()
	if err != nil {
		return nil, err
	}

	nativeSettings, _, err = settings.GetSettings()
	if err != nil {
		return nil, err
	}

	return nativeSettings, nil
}

func (s *GitlabSettings) GetFqdn() (fqdn string, err error) {
	gitlab, err := s.GetGitlab()
	if err != nil {
		return "", err
	}

	fqdn, err = gitlab.GetFqdn()
	if err != nil {
		return "", err
	}

	return fqdn, nil
}

func (s *GitlabSettings) GetGitlab() (gitlab *GitlabInstance, err error) {
	if s.gitlab == nil {
		return nil, tracederrors.TracedError("gitlab not set")
	}

	return s.gitlab, nil
}

func (s *GitlabSettings) GetNativeClient() (nativeClient *gitlab.Client, err error) {
	gitlab, err := s.GetGitlab()
	if err != nil {
		return nil, err
	}

	nativeClient, err = gitlab.GetNativeClient()
	if err != nil {
		return nil, err
	}

	return nativeClient, nil
}

func (s *GitlabSettings) GetNativeSettingsService() (nativeSettingsService *gitlab.SettingsService, err error) {
	nativeClient, err := s.GetNativeClient()
	if err != nil {
		return nil, err
	}

	nativeSettingsService = nativeClient.Settings
	return nativeSettingsService, nil
}

func (s *GitlabSettings) IsAutoDevopsEnabled() (isAutoDevopsEnabled bool, err error) {
	nativeSettings, err := s.GetCurrentSettingsNative()
	if err != nil {
		return false, err
	}

	isAutoDevopsEnabled = nativeSettings.AutoDevOpsEnabled

	return isAutoDevopsEnabled, nil
}

func (s *GitlabSettings) IsSignupEnabled() (isSignupEnabled bool, err error) {
	nativeSettings, err := s.GetCurrentSettingsNative()
	if err != nil {
		return false, err
	}

	isSignupEnabled = nativeSettings.SignupEnabled

	return isSignupEnabled, nil
}

func (s *GitlabSettings) MustDisableAutoDevos(verbose bool) {
	err := s.DisableAutoDevops(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (s *GitlabSettings) SetGitlab(gitlab *GitlabInstance) (err error) {
	if gitlab == nil {
		return tracederrors.TracedError("gitlab is nil")
	}

	s.gitlab = gitlab

	return nil
}
