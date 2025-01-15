package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type AuthenticationOptionsHandlerService struct{}

func AuthenticationOptionsHandler() (a *AuthenticationOptionsHandlerService) {
	return NewAuthenticationOptionsHandlerService()
}

func NewAuthenticationOptionsHandlerService() (a *AuthenticationOptionsHandlerService) {
	return new(AuthenticationOptionsHandlerService)
}

func (a *AuthenticationOptionsHandlerService) GetAuthenticationoptionsForService(authentiationOptions []AuthenticationOption, serviceName string) (authOption AuthenticationOption, err error) {
	if serviceName == "" {
		return nil, tracederrors.TracedErrorEmptyString("serviceName")
	}

	for _, authOption = range authentiationOptions {
		isAuthenticating, err := authOption.IsAuthenticatingAgainst(serviceName)
		if err != nil {
			return nil, err
		}

		if isAuthenticating {
			return authOption, nil
		}
	}

	return nil, tracederrors.TracedErrorf(
		"No authenticationOptions for '%s' found. Checked '%d' authenticationOptions in total.",
		serviceName,
		len(authentiationOptions),
	)
}

func (a *AuthenticationOptionsHandlerService) GetAuthenticationoptionsForServiceByUrl(authenticationOptions []AuthenticationOption, url *URL) (authOption AuthenticationOption, err error) {
	if url == nil {
		return nil, tracederrors.TracedErrorNil("url")
	}

	urlString, err := url.GetAsString()
	if err != nil {
		return nil, err
	}

	authOption, err = a.GetAuthenticationoptionsForService(authenticationOptions, urlString)
	if err != nil {
		return nil, err
	}

	return authOption, nil
}

func (a *AuthenticationOptionsHandlerService) MustGetAuthenticationoptionsForService(authentiationOptions []AuthenticationOption, serviceName string) (authOption AuthenticationOption) {
	authOption, err := a.GetAuthenticationoptionsForService(authentiationOptions, serviceName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return authOption
}

func (a *AuthenticationOptionsHandlerService) MustGetAuthenticationoptionsForServiceByUrl(authenticationOptions []AuthenticationOption, url *URL) (authOption AuthenticationOption) {
	authOption, err := a.GetAuthenticationoptionsForServiceByUrl(authenticationOptions, url)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return authOption
}
