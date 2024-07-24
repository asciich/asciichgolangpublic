package asciichgolangpublic

type AuthenticationOption interface {
	IsAuthenticatingAgainst(serviceName string) (isAuthenticatingAgainst bool, err error)
	IsVerbose() (isVerbose bool)
}
