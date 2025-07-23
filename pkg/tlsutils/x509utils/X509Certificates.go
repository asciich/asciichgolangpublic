package x509utils

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/asciich/asciichgolangpublic/files"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/shellutils/shelllinehandler"
	"github.com/asciich/asciichgolangpublic/tempfiles"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

func GetDefaultHandler() (certHandler X509CertificateHandler) {
	return GetNativeX509CertificateHandler()
}

func CreateRootCa(ctx context.Context, options *X509CreateCertificateOptions) (*X509CertKeyPair, error) {
	return GetDefaultHandler().CreateRootCaCertificate(ctx, options)
}

func CreateSignedIntermediateCertificate(ctx context.Context, options *X509CreateCertificateOptions, rootCaCertAndKey *X509CertKeyPair) (*X509CertKeyPair, error) {
	return GetDefaultHandler().CreateSignedIntermediateCertificate(ctx, options, rootCaCertAndKey)
}

func CreateSignedEndEndityCertificate(ctx context.Context, options *X509CreateCertificateOptions, caCertAndKey *X509CertKeyPair) (endEndityCertAndKey *X509CertKeyPair, err error) {
	return GetDefaultHandler().CreateSignedEndEndityCertificate(ctx, options, caCertAndKey)
}

// ================================
// TODO rewrite/ remove from here:
// ================================
type X509CertificatesService struct {
}

// Deprecated: Is reimplemented without additional X509CertificatesService struct.
func NewX509CertificatesService() (x *X509CertificatesService) {
	return new(X509CertificatesService)
}

// Deprecated: Is reimplemented without additional X509CertificatesService struct.
func X509Certificates() (x509Certificaets *X509CertificatesService) {
	return new(X509CertificatesService)
}

func (c *X509CertificatesService) CreateIntermediateCertificateIntoDirectory(ctx context.Context, createOptions *X509CreateCertificateOptions) (directoryContianingCreatedCertAndKey files.Directory, err error) {
	if createOptions == nil {
		return nil, tracederrors.TracedError("createOptions is nil")
	}

	if !createOptions.GetUseTemporaryDirectory() {
		return nil, tracederrors.TracedError("Only implemented for temporary directory")
	}

	directoryToUse, err := tempfiles.CreateEmptyTemporaryDirectory(true)
	if err != nil {
		return nil, err
	}

	directoryPathToUse, err := directoryToUse.GetLocalPath()
	if err != nil {
		return nil, err
	}

	subjectString, err := createOptions.GetSubjectStringForOpenssl()
	if err != nil {
		return nil, err
	}

	logging.LogInfof("Going to create new intermediate certificate for '%v'", subjectString)

	sslCommand := []string{
		"openssl",
		"genrsa",
		"-out",
		filepath.Join(directoryPathToUse, "intermediateCertificate.key"),
		"4096",
	}

	_, err = commandexecutor.Bash().RunCommand(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: sslCommand,
		},
	)
	if err != nil {
		return nil, err
	}

	logging.LogInfof("Created intermediate certificate in temporary directory: '%v'", directoryPathToUse)

	return directoryToUse, nil
}

func (c *X509CertificatesService) CreateRootCaIntoDirectory(ctx context.Context, createOptions *X509CreateCertificateOptions) (directoryContianingCreatedCertAndKey files.Directory, err error) {
	if createOptions == nil {
		return nil, tracederrors.TracedError("createOptions is nil")
	}

	if !createOptions.GetUseTemporaryDirectory() {
		return nil, tracederrors.TracedError("Only implemented for temporary directory")
	}

	directoryToUse, err := tempfiles.CreateEmptyTemporaryDirectory(true)
	if err != nil {
		return nil, err
	}

	directoryPathToUse, err := directoryToUse.GetLocalPath()
	if err != nil {
		return nil, err
	}

	subjectString, err := createOptions.GetSubjectStringForOpenssl()
	if err != nil {
		return nil, err
	}

	logging.LogInfof("Going to create new RootCA for '%v'", subjectString)

	sslCommand := []string{
		"openssl",
		"req",
		"-x509",
		"-sha256",
		"-days",
		"356",
		"-nodes",
		"-newkey",
		"rsa:4096",
		"-subj",
		subjectString,
		"-keyout",
		"rootCA.key",
		"-out",
		"rootCA.crt",
	}

	joinedSslCommand, err := shelllinehandler.Join(sslCommand)
	if err != nil {
		return nil, err
	}

	createCommand := []string{
		"bash",
		"-c",
		fmt.Sprintf("cd '%v' && %v", directoryPathToUse, joinedSslCommand),
	}

	_, err = commandexecutor.Bash().RunCommand(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: createCommand,
		},
	)
	if err != nil {
		return nil, err
	}

	logging.LogInfof("Created root ca in temporary directory: '%v'", directoryPathToUse)

	return directoryToUse, nil
}

func (c *X509CertificatesService) CreateSignedCertificate(createOptions *X509CreateCertificateOptions) (err error) {
	return tracederrors.TracedErrorNotImplemented()
	/* TODO move gopass part
	if createOptions == nil {
		return tracederrors.TracedError("createOptions is nil")
	}

	outputKeyOnStdout := false

	keyPath, err := createOptions.GetKeyOutputFilePath()
	if err != nil {
		return err
	}

	if keyPath == "-" {
		outputKeyOnStdout = true
		tempKeyFile, err := tempfiles.CreateEmptyTemporaryFile(createOptions.Verbose)
		if err != nil {
			return err
		}

		defer tempKeyFile.SecurelyDelete(createOptions.Verbose)

		keyPath, err = tempKeyFile.GetLocalPath()
		if err != nil {
			return err
		}
	}

	var certPath string = ""
	if createOptions.IsCertificateOutputFilePathSet() {
		certPath, err = createOptions.GetCertificateOutputFilePath()
		if err != nil {
			return err
		}
	} else {
		tempCertFile, err := tempfiles.CreateEmptyTemporaryFile(createOptions.Verbose)
		if err != nil {
			return err
		}

		defer tempCertFile.SecurelyDelete(createOptions.Verbose)

		certPath, err = tempCertFile.GetLocalPath()
		if err != nil {
			return err
		}
	}

	csrPath, err := tempfiles.CreateEmptyTemporaryFileAndGetPath(createOptions.Verbose)
	if err != nil {
		return err
	}

	commonName, err := createOptions.GetCommonName()
	if err != nil {
		return err
	}

	subjectString, err := createOptions.GetSubjectStringForOpenssl()
	if err != nil {
		return err
	}

	createKeyCommand := []string{
		"openssl",
		"req",
		"-nodes",
		"-newkey",
		"rsa:4096",
		"-keyout",
		keyPath,
		"-out",
		csrPath,
		"-subj",
		subjectString,
		"-addext",
		fmt.Sprintf("subjectAltName = DNS:%s", commonName),
	}

	_, err = commandexecutor.Bash().RunCommand(
		&parameteroptions.RunCommandOptions{
			Command: createKeyCommand,
			Verbose: createOptions.Verbose,
		},
	)
	if err != nil {
		return err
	}

	if createOptions.Verbose {
		logging.LogInfof("Created key file '%v' and CSR '%v' for '%v'", keyPath, csrPath, commonName)
	}

	if createOptions.Verbose {
		opensslCsrInfoCmd := []string{
			"openssl",
			"req",
			"-noout",
			"-text",
			"-in",
			csrPath,
		}

		csrInfo, err := commandexecutor.Bash().RunCommandAndGetStdoutAsString(
			&parameteroptions.RunCommandOptions{
				Command: opensslCsrInfoCmd,
			},
		)
		if err != nil {
			return err
		}

		logging.LogInfof("Created CSR info:\n'%s'", csrInfo)
	}

	signingKey, err := createOptions.GetIntermediateCertificateKeyGopassCredential()
	if err != nil {
		return err
	}

	signingKeyFile, err := signingKey.WriteIntoTemporaryFile(createOptions.Verbose)
	if err != nil {
		return err
	}
	defer signingKeyFile.SecurelyDelete(createOptions.Verbose)

	signingKeyFilePath, err := signingKeyFile.GetLocalPath()
	if err != nil {
		return err
	}

	signingCert, err := createOptions.GetIntermediateCertificateGopassCredential()
	if err != nil {
		return err
	}

	signingCertFile, err := signingCert.WriteIntoTemporaryFile(createOptions.Verbose)
	if err != nil {
		return err
	}
	defer signingCertFile.SecurelyDelete(createOptions.Verbose)

	signingCertFilePath, err := signingCertFile.GetLocalPath()
	if err != nil {
		return err
	}

	signingConfig := ""
	signingConfig += "[req]\n"
	signingConfig += "req_extensions = req_ext\n"
	signingConfig += "prompt = no\n"
	signingConfig += "\n"
	signingConfig += "[req_ext]\n"
	signingConfig += "subjectAltName = @alt_names\n"
	signingConfig += "\n"
	signingConfig += "[alt_names]\n"
	signingConfig += "DNS.1 = " + commonName + "\n"
	for i, san := range createOptions.AdditionalSans {
		signingConfig += fmt.Sprintf("DNS.%d = %s\n", i+2, san)
		if createOptions.Verbose {
			logging.LogInfof("Added additional SAN '%s' for commonName '%s'.", san, commonName)
		}
	}
	if createOptions.Verbose {
		logging.LogInfof("Added '%d' SAN's to sing with '%s'.", len(createOptions.AdditionalSans), commonName)
	}

	signingConfigFile, err := tempfiles.CreateFromString(signingConfig, createOptions.Verbose)
	if err != nil {
		return err
	}
	defer signingConfigFile.SecurelyDelete(createOptions.Verbose)

	signingConfigFilePath, err := signingConfigFile.GetLocalPath()
	if err != nil {
		return err
	}

	serial, err := c.GetNextCaSerialNumberAsStringFromGopass(createOptions.Verbose)
	if err != nil {
		return err
	}

	signCommand := []string{
		"openssl",
		"x509",
		"-req",
		"-days",
		"45",
		"-in",
		csrPath,
		"-CA",
		signingCertFilePath,
		"-CAkey",
		signingKeyFilePath,
		"-set_serial",
		serial,
		"-out",
		certPath,
		"-extensions",
		"req_ext",
		"-extfile",
		signingConfigFilePath,
	}
	_, err = commandexecutor.Bash().RunCommand(
		&parameteroptions.RunCommandOptions{
			Command: signCommand,
			Verbose: createOptions.Verbose,
		},
	)
	if err != nil {
		return err
	}

	if createOptions.Verbose {
		logging.LogInfof("Created certificate file: '%v'", certPath)
	}

	if createOptions.Verbose {
		certInfoCommand := []string{
			"openssl",
			"x509",
			"-noout",
			"-text",
			"-in",
			certPath,
		}

		certificateInfo, err := commandexecutor.Bash().RunCommandAndGetStdoutAsString(
			&parameteroptions.RunCommandOptions{
				Command: certInfoCommand,
			},
		)
		if err != nil {
			return err
		}

		logging.LogInfof("Created certificate info:\n%s", certificateInfo)
	}

	certFile, err := files.GetLocalFileByPath(certPath)
	if err != nil {
		return err
	}

	err = Gopass().InsertFile(certFile, &GopassSecretOptions{
		SecretRootDirectoryPath: "internal_ca/created_certificates/" + commonName,
		SecretBasename:          commonName + ".crt",
		Verbose:                 createOptions.Verbose,
		Overwrite:               createOptions.OverwriteExistingCertificateInGopass,
	})
	if err != nil {
		return err
	}

	if outputKeyOnStdout {
		keyFile, err := files.GetLocalFileByPath(keyPath)
		if err != nil {
			return err
		}

		err = keyFile.PrintContentOnStdout()
		if err != nil {
			return err
		}
	}

	return nil
	*/
}

func (c *X509CertificatesService) CreateSignedIntermediateCertificateAndAddToGopass(createOptions *X509CreateCertificateOptions, rootCaInGopass *parameteroptions.GopassSecretOptions, intermediateGopassOptions *parameteroptions.GopassSecretOptions) (err error) {
	return tracederrors.TracedErrorNotImplemented()
	/* TODO move to gopass
	if createOptions == nil {
		return tracederrors.TracedError("createOptions is nil")
	}

	if rootCaInGopass == nil {
		return tracederrors.TracedError("rootCaInGopass is nil")
	}

	if intermediateGopassOptions == nil {
		return tracederrors.TracedError("intermediateGopassOptions is nil")
	}

	rootCertOptions := rootCaInGopass.GetDeepCopy()
	rootCertOptions.SecretBasename = "rootCa.crt"
	rootCertFile, err := Gopass().WriteSecretIntoTemporaryFile(rootCertOptions)
	if err != nil {
		return err
	}
	defer rootCertFile.SecurelyDelete(createOptions.Verbose)

	rootKeyOptions := rootCaInGopass.GetDeepCopy()
	rootKeyOptions.SecretBasename = "rootCa.key"
	rootKeyFile, err := Gopass().WriteSecretIntoTemporaryFile(rootKeyOptions)
	if err != nil {
		return err
	}
	defer rootKeyFile.SecurelyDelete(createOptions.Verbose)

	createOptionsToUse := createOptions.GetDeepCopy()
	createOptionsToUse.UseTemporaryDirectory = true
	intermediateDirectory, err := c.CreateIntermediateCertificateIntoDirectory(createOptionsToUse)
	if err != nil {
		return err
	}

	intermediateCertFile, err := intermediateDirectory.GetFileInDirectory("intermediateCertificate.crt")
	if err != nil {
		return err
	}
	defer intermediateCertFile.SecurelyDelete(createOptions.Verbose)

	intermediateKeyFile, err := intermediateDirectory.GetFileInDirectory("intermediateCertificate.key")
	if err != nil {
		return err
	}
	defer intermediateKeyFile.SecurelyDelete(createOptions.Verbose)

	signingRequestFile, err := tempfiles.CreateEmptyTemporaryFile(createOptions.Verbose)
	if err != nil {
		return err
	}
	defer signingRequestFile.SecurelyDelete(createOptions.Verbose)

	signingOptions := NewX509SignCertificateOptions()
	signingOptions.CertFileUsedForSigning = rootCertFile
	signingOptions.KeyFileUsedForSigning = rootKeyFile
	signingOptions.KeyFileToSign = intermediateKeyFile
	signingOptions.OutputCertificateFile = intermediateCertFile
	signingOptions.SigningRequestFile = signingRequestFile
	signingOptions.CommonName = createOptions.CommonName
	signingOptions.CountryName = createOptions.CountryName
	signingOptions.Locality = createOptions.Locality
	signingOptions.Verbose = createOptions.Verbose

	err = c.SignIntermediateCertificate(signingOptions)
	if err != nil {
		return err
	}

	gopassInsertOptions := intermediateGopassOptions.GetDeepCopy()
	gopassInsertOptions.SecretBasename = "intermediateCertificate.key"
	err = Gopass().InsertFile(intermediateKeyFile, gopassInsertOptions)
	if err != nil {
		return err
	}

	gopassInsertOptions = intermediateGopassOptions.GetDeepCopy()
	gopassInsertOptions.SecretBasename = "intermediateCertificate.crt"
	err = Gopass().InsertFile(intermediateCertFile, gopassInsertOptions)
	if err != nil {
		return err
	}

	return nil
	*/
}

func (c *X509CertificatesService) CreateSigningRequestFile(signOptions *X509SignCertificateOptions) (err error) {
	if signOptions == nil {
		return tracederrors.TracedError("signOptions is nil")
	}

	keyFileToSignPath, err := signOptions.GetKeyFileToSignPath()
	if err != nil {
		return err
	}

	signingRequestFilePath, err := signOptions.GetSigningRequestFilePath()
	if err != nil {
		return err
	}

	subjectToSign, err := signOptions.GetSubjectToSign()
	if err != nil {
		return err
	}

	openSslConfigFile, err := tempfiles.CreateEmptyTemporaryFile(signOptions.Verbose)
	if err != nil {
		return err
	}

	const extensionName = "v3_req"

	err = openSslConfigFile.WriteString(
		" [ req ]\n"+
			"req_extensions = "+extensionName+"\n"+
			"x509_extensions = "+extensionName+"\n"+
			"\n"+
			"[ "+extensionName+" ]\n"+
			"basicConstraints = CA:TRUE\n",
		signOptions.Verbose,
	)
	if err != nil {
		return err
	}

	openSslConfigFilePath, err := openSslConfigFile.GetLocalPath()
	if err != nil {
		return err
	}

	if signOptions.Verbose {
		logging.LogInfof("Generated openssl configuration for signing request: '%v'.", openSslConfigFilePath)
	}

	signCommand := []string{
		"openssl",
		"req",
		"-new",
		"-sha256",
		"-config",
		openSslConfigFilePath,
		"-extensions",
		extensionName,
		"-subj",
		subjectToSign,
		"-key",
		keyFileToSignPath,
		"-out",
		signingRequestFilePath,
	}

	_, err = commandexecutor.Bash().RunCommand(
		contextutils.GetVerbosityContextByBool(signOptions.Verbose),
		&parameteroptions.RunCommandOptions{
			Command: signCommand,
		},
	)
	if err != nil {
		return err
	}

	if signOptions.Verbose {
		logging.LogInfof("Created signing request file '%v' for '%v' with key located at '%v'.", signingRequestFilePath, subjectToSign, keyFileToSignPath)
	}

	return nil
}

func (c *X509CertificatesService) GetNextCaSerialNumberAsStringFromGopass(verbose bool) (serial string, err error) {
	return "", tracederrors.TracedErrorNotImplemented()
	/* TODO move to gopass
	nextFreeNumberFromEnvVar := os.Getenv("OVERRIDE_NEXT_CA_SERIAL_NUMBER_FROM_GOPASS")
	if len(nextFreeNumberFromEnvVar) > 0 {
		if verbose {
			logging.LogInfof("GetNextCaSerialNumberAsStringFromGopass: OVERRIDE_NEXT_CA_SERIAL_NUMBER_FROM_GOPASS is set to '%s'", nextFreeNumberFromEnvVar)
		}

		return nextFreeNumberFromEnvVar, nil
	}

	serialCredential, err := Gopass().GetCredential(
		&GopassSecretOptions{
			SecretRootDirectoryPath: "internal_ca/root_ca",
			SecretBasename:          "serial_counter",
		},
	)
	if err != nil {
		return "", err
	}

	err = serialCredential.IncrementIntValue()
	if err != nil {
		return "", err
	}

	serial, err = serialCredential.GetAsString()
	if err != nil {
		return "", err
	}

	serial = strings.TrimSpace(serial)

	return serial, nil
	*/
}

/*
func (c *X509CertificatesService) IsCertificateFileSignedByCertificateFile(thisCertificateFile *X509CertificateFile, isSignedByThisCertificateFile files.File, verbose bool) (isSignedBy bool, err error) {
	if thisCertificateFile == nil {
		return false, tracederrors.TracedError("thisCertificateFile is nil")
	}

	if isSignedByThisCertificateFile == nil {
		return false, tracederrors.TracedError("isSignedByThisCertificateFile is nil")
	}

	toCheckCert, err := thisCertificateFile.GetAsX509Certificate()
	if err != nil {
		return false, err
	}

	isSignedBy, err = toCheckCert.IsSignedByCertificateFile(isSignedByThisCertificateFile, verbose)
	if err != nil {
		return false, err
	}

	return isSignedBy, err
}
*/

func (c *X509CertificatesService) SignIntermediateCertificate(signOptions *X509SignCertificateOptions) (err error) {
	if signOptions == nil {
		return tracederrors.TracedError("signOptions is nil")
	}

	keyFileToUseForSigning, err := signOptions.GetKeyFileUsedForSigning()
	if err != nil {
		return err
	}

	keyFileToUseForSigningPath, err := keyFileToUseForSigning.GetLocalPath()
	if err != nil {
		return err
	}

	certFileToUseForSigning, err := signOptions.GetCertFileUsedForSigning()
	if err != nil {
		return err
	}

	certFileToUseForSigningPath, err := certFileToUseForSigning.GetLocalPath()
	if err != nil {
		return err
	}

	outputCertificateFile, err := signOptions.GetOutputCertificateFile()
	if err != nil {
		return err
	}

	outputCertificateFilePath, err := outputCertificateFile.GetLocalPath()
	if err != nil {
		return err
	}

	singingRequestFile, err := tempfiles.CreateEmptyTemporaryFile(signOptions.Verbose)
	if err != nil {
		return err
	}

	signOptionsToUse := signOptions.GetDeepCopy()
	signOptionsToUse.SigningRequestFile = singingRequestFile
	err = c.CreateSigningRequestFile(signOptionsToUse)
	if err != nil {
		return err
	}

	signingRequestFilePath, err := singingRequestFile.GetLocalPath()
	if err != nil {
		return err
	}

	openSslConfigFile, err := tempfiles.CreateEmptyTemporaryFile(signOptions.Verbose)
	if err != nil {
		return err
	}

	const extensionName = "v3_req"

	err = openSslConfigFile.WriteString(
		" [ req ]\n"+
			"req_extensions = "+extensionName+"\n"+
			"x509_extensions = "+extensionName+"\n"+
			"\n"+
			"[ "+extensionName+" ]\n"+
			"basicConstraints = CA:TRUE\n",
		signOptions.Verbose,
	)
	if err != nil {
		return err
	}

	openSslConfigFilePath, err := openSslConfigFile.GetLocalPath()
	if err != nil {
		return err
	}

	if signOptions.Verbose {
		logging.LogInfof("Generated openssl configuration for signing: '%v'.", openSslConfigFilePath)
	}

	serial, err := c.GetNextCaSerialNumberAsStringFromGopass(signOptions.Verbose)
	if err != nil {
		return err
	}

	signCommand := []string{
		"openssl",
		"x509",
		"-req",
		"-days",
		"90",
		"-extfile",
		openSslConfigFilePath,
		"-extensions",
		extensionName,
		"-in",
		signingRequestFilePath,
		"-CA",
		certFileToUseForSigningPath,
		"-CAkey",
		keyFileToUseForSigningPath,
		"-set_serial",
		serial,
		"-out",
		outputCertificateFilePath,
	}

	_, err = commandexecutor.Bash().RunCommand(
		contextutils.GetVerbosityContextByBool(signOptions.Verbose),
		&parameteroptions.RunCommandOptions{
			Command: signCommand,
		},
	)
	if err != nil {
		return err
	}

	if signOptions.Verbose {
		logging.LogInfof("Signed intermediate certificate. Certificate stored as '%v'", outputCertificateFilePath)
	}

	return nil
}

func (x *X509CertificatesService) MustCreateSigningRequestFile(signOptions *X509SignCertificateOptions) {
	err := x.CreateSigningRequestFile(signOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}
