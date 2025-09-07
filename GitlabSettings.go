package asciichgolangpublic

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type GitlabSettings struct {
	gitlab *GitlabInstance
}

func NewGitlabSettings() (gitlabSettings *GitlabSettings) {
	return new(GitlabSettings)
}

func (s *GitlabSettings) DisableAutoDevops(ctx context.Context) (err error) {
	isAutoDevopsEnabled, err := s.IsAutoDevopsEnabled()
	if err != nil {
		return err
	}

	fqdn, err := s.GetFqdn()
	if err != nil {
		return err
	}

	if isAutoDevopsEnabled {
		logging.LogInfoByCtxf(ctx, "Disable AutoDevops on gitlab '%s' started.", fqdn)

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
		logging.LogChangedByCtxf(ctx, "Dislabed AutoDevops on gitlab '%s'.", fqdn)

	} else {
		logging.LogInfoByCtxf(ctx, "Autodevops on gitlab '%s' already disabled.", fqdn)
	}

	logging.LogInfoByCtxf(ctx, "Disable AutoDevops on gitlab '%s' finished.", fqdn)

	return nil
}

func (s *GitlabSettings) DisableSignup(ctx context.Context) (err error) {
	isSignupEnabled, err := s.IsSignupEnabled()
	if err != nil {
		return err
	}

	fqdn, err := s.GetFqdn()
	if err != nil {
		return err
	}

	if isSignupEnabled {
		logging.LogInfoByCtxf(ctx, "Disable signup to gitlab '%s' started.", fqdn)

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
		logging.LogChangedByCtxf(ctx, "Dislabed signup to gitlab '%s'.", fqdn)

	} else {
		logging.LogInfoByCtxf(ctx, "Signup to gitlab '%s' already disabled.", fqdn)
	}
	logging.LogInfoByCtxf(ctx, "Disable signup to gitlab '%s' finished.", fqdn)
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

func (s *GitlabSettings) SetGitlab(gitlab *GitlabInstance) (err error) {
	if gitlab == nil {
		return tracederrors.TracedError("gitlab is nil")
	}

	s.gitlab = gitlab

	return nil
}
