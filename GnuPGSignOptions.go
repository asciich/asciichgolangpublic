package asciichgolangpublic

type GnuPGSignOptions struct {
	Verbose      bool
	DetachedSign bool
	AsciiArmor   bool
}

func NewGnuPGSignOptions() (g *GnuPGSignOptions) {
	return new(GnuPGSignOptions)
}

func (g *GnuPGSignOptions) GetAsciiArmor() (asciiArmor bool, err error) {

	return g.AsciiArmor, nil
}

func (g *GnuPGSignOptions) GetDetachedSign() (detachedSign bool, err error) {

	return g.DetachedSign, nil
}

func (g *GnuPGSignOptions) GetVerbose() (verbose bool, err error) {

	return g.Verbose, nil
}

func (g *GnuPGSignOptions) MustGetAsciiArmor() (asciiArmor bool) {
	asciiArmor, err := g.GetAsciiArmor()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return asciiArmor
}

func (g *GnuPGSignOptions) MustGetDetachedSign() (detachedSign bool) {
	detachedSign, err := g.GetDetachedSign()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return detachedSign
}

func (g *GnuPGSignOptions) MustGetVerbose() (verbose bool) {
	verbose, err := g.GetVerbose()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return verbose
}

func (g *GnuPGSignOptions) MustSetAsciiArmor(asciiArmor bool) {
	err := g.SetAsciiArmor(asciiArmor)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GnuPGSignOptions) MustSetDetachedSign(detachedSign bool) {
	err := g.SetDetachedSign(detachedSign)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GnuPGSignOptions) MustSetVerbose(verbose bool) {
	err := g.SetVerbose(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GnuPGSignOptions) SetAsciiArmor(asciiArmor bool) (err error) {
	g.AsciiArmor = asciiArmor

	return nil
}

func (g *GnuPGSignOptions) SetDetachedSign(detachedSign bool) (err error) {
	g.DetachedSign = detachedSign

	return nil
}

func (g *GnuPGSignOptions) SetVerbose(verbose bool) (err error) {
	g.Verbose = verbose

	return nil
}
