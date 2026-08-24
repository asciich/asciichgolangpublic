// Package commandexecutorx509utils provides command-based implementation for X509 certificate operations.
// This implementation uses the command executor for file I/O (read/write) on remote machines,
// while certificate generation and signing is handled by genericx509utils using native Go crypto.
package commandexecutorx509utils

import (
	"context"
	"crypto"
	"crypto/x509"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils/genericx509utils"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils/x509options"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func ReadCertificateFromFile(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, pathToRead string) (cert *x509.Certificate, err error) {
	if commandExecutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExecutor")
	}

	if pathToRead == "" {
		return nil, tracederrors.TracedErrorEmptyString("pathToRead")
	}

	logging.LogInfoByCtxf(ctx, "Read certificate from file '%s' using command executor started.", pathToRead)

	content, err := commandexecutorfile.ReadAsBytes(commandExecutor, pathToRead)
	if err != nil {
		return nil, err
	}

	cert, err = genericx509utils.ReadCertFromBytes(content)
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Read certificate from file '%s' using command executor finished.", pathToRead)

	return cert, nil
}

func WriteCertificateToFile(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, cert *x509.Certificate, pathToWrite string) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if cert == nil {
		return tracederrors.TracedErrorNil("cert")
	}

	if pathToWrite == "" {
		return tracederrors.TracedErrorEmptyString("pathToWrite")
	}

	logging.LogInfoByCtxf(ctx, "Write certificate to file '%s' using command executor started.", pathToWrite)

	certBytes, err := genericx509utils.WriteCertAsBytes(cert)
	if err != nil {
		return err
	}

	err = commandexecutorfile.WriteBytes(ctx, commandExecutor, pathToWrite, certBytes, &filesoptions.WriteOptions{})
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Write certificate to file '%s' using command executor finished.", pathToWrite)

	return nil
}

func GeneratePrivateKey(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) (privateKey crypto.PrivateKey, err error) {
	if commandExecutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExecutor")
	}

	return genericx509utils.GeneratePrivateKey(ctx)
}

func CreateRootCaCertificate(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, options *x509options.X509CreateCertificateOptions) (caCertAndKey *genericx509utils.X509CertKeyPair, err error) {
	if commandExecutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExecutor")
	}

	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	return genericx509utils.CreateRootCaCertificate(ctx, options)
}

func CreateIntermediateCertificate(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, options *x509options.X509CreateCertificateOptions) (intermediateCert *genericx509utils.X509CertKeyPair, err error) {
	if commandExecutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExecutor")
	}

	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	return genericx509utils.CreateIntermediateCertificate(ctx, options)
}

func CreateSelfSignedCertificate(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, options *x509options.X509CreateCertificateOptions) (selfSignedCertAndKey *genericx509utils.X509CertKeyPair, err error) {
	if commandExecutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExecutor")
	}

	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	return genericx509utils.CreateSelfSignedCertificate(ctx, options)
}

func CreateSignedIntermediateCertificate(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, options *x509options.X509CreateCertificateOptions, rootCaCertAndKey *genericx509utils.X509CertKeyPair) (intermediateCertAndKey *genericx509utils.X509CertKeyPair, err error) {
	if commandExecutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExecutor")
	}

	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	if rootCaCertAndKey == nil {
		return nil, tracederrors.TracedErrorNil("rootCaCertAndKey")
	}

	return genericx509utils.CreateSignedIntermediateCertificate(ctx, options, rootCaCertAndKey)
}

func CreateSignedEndEntityCertificate(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, options *x509options.X509CreateCertificateOptions, caCertAndKey *genericx509utils.X509CertKeyPair) (endEntityCertAndKey *genericx509utils.X509CertKeyPair, err error) {
	if commandExecutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExecutor")
	}

	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	if caCertAndKey == nil {
		return nil, tracederrors.TracedErrorNil("caCertAndKey")
	}

	return genericx509utils.CreateSignedEndEntityCertificate(ctx, options, caCertAndKey)
}
