package authenticationoptions

type AuthenticationOption interface {
	IsAuthenticatingAgainst(serviceName string) (isAuthenticatingAgainst bool, err error)
}
