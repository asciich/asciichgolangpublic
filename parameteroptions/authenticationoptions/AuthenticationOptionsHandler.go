package authenticationoptions

import (
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/urlsutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/tracederrors"
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

func (a *AuthenticationOptionsHandlerService) GetAuthenticationoptionsForServiceByUrl(authenticationOptions []AuthenticationOption, url *urlsutils.URL) (authOption AuthenticationOption, err error) {
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
