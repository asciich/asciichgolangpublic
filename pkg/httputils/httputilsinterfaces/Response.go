package httputilsinterfaces

import (
	"context"
	"crypto/x509"
)

type Response interface {
	CheckStatusCode(expectedStatusCodes []int) error
	GetBodyAsBytes() (body []byte, err error)
	GetBodyAsString() (body string, err error)
	IsStatusCode(expectedStatusCode int) bool
	IsStatusCode200Ok() bool
	SetBody(body []byte) (err error)
	SetStatusCode(statusCode int) (err error)
	RunJqQueryAgainstBody(query string) (result string, err error)
	RunYqQueryAgainstBody(query string) (result string, err error)

	// GetServerEndEntitiyCertificate returns the server's end entity (leaf) certificate
	GetServerEndEntitiyCertificate(ctx context.Context) (*x509.Certificate, error)

	// GetServerCertificateChain returns the full certificate chain from the server
	GetServerCertificateChain(ctx context.Context) ([]*x509.Certificate, error)

	// LogCertInfo logs the server's certificate information
	LogCertInfo(ctx context.Context) error
}
