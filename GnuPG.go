package asciichgolangpublic

type GnuPGService struct {
}

func GnuPG() (gnuPG *GnuPGService) {
	return new(GnuPGService)
}

func NewGnuPGService() (g *GnuPGService) {
	return new(GnuPGService)
}

func (g *GnuPGService) CheckSignatureValid(signatureFile File, verbose bool) (err error) {
	if signatureFile == nil {
		return TracedErrorNil("signatureFile")
	}

	err = signatureFile.CheckIsLocalFile(verbose)
	if err != nil {
		return TracedErrorf("Only implemented for local available files: %w", err)
	}

	path, hostDescription, err := signatureFile.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	if verbose {
		LogInfof(
			"Validate GnuPG signature from '%s' on host '%s' started.",
			path,
			hostDescription,
		)
	}

	_, err = Bash().RunCommand(
		&RunCommandOptions{
			Command: []string{"gpg", "--verify", path},
			Verbose: verbose,
		},
	)

	if verbose {
		LogInfof(
			"GnuPG signature from '%s' on host '%s' validated.",
			path,
			hostDescription,
		)
	}

	return nil
}

func (g *GnuPGService) MustCheckSignatureValid(signatureFile File, verbose bool) {
	err := g.CheckSignatureValid(signatureFile, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GnuPGService) MustSignFile(fileToSign File, options *GnuPGSignOptions) {
	err := g.SignFile(fileToSign, options)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GnuPGService) SignFile(fileToSign File, options *GnuPGSignOptions) (err error) {
	if fileToSign == nil {
		return TracedError("fileToSign is nil")
	}

	if options == nil {
		return TracedError("options is nil")
	}

	err = fileToSign.CheckIsLocalFile(options.Verbose)
	if err != nil {
		return TracedErrorf("Only implemented for local available files: %w", err)
	}

	path, err := fileToSign.GetPath()
	if err != nil {
		return err
	}

	if options.Verbose {
		LogInfof("Sign '%s' using gnupg started.", path)
	}

	if !options.AsciiArmor {
		return TracedError("Only implemented for asciiArmor at the moment")
	}

	if !options.DetachedSign {
		return TracedError("Only implemented for DetachedSign at the moment")
	}

	signaturePath := path + ".asc"
	signatureFile, err := GetLocalFileByPath(signaturePath)
	if err != nil {
		return err
	}

	if err = signatureFile.Delete(options.Verbose); err != nil {
		return err
	}

	signCommand := []string{
		"gpg",
		"--armor",
		"--detach-sig",
		path,
	}

	_, err = Bash().RunCommand(
		&RunCommandOptions{
			Command: signCommand,
			Verbose: options.Verbose,
		},
	)
	if err != nil {
		return err
	}

	signatureFileExists, err := signatureFile.Exists(false)
	if err != nil {
		return err
	}

	if !signatureFileExists {
		return TracedErrorf(
			"Signing '%s' failed. Expected signature file '%s' does not exits.",
			path,
			signaturePath,
		)
	}

	if options.Verbose {
		LogInfof("Sign '%s' using gnupg finished.", path)
	}

	return nil
}
