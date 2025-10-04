package gnupg

import "github.com/asciich/asciichgolangpublic/pkg/logging"

type GnuPGSignOptions struct {
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

func (g *GnuPGSignOptions) MustGetAsciiArmor() (asciiArmor bool) {
	asciiArmor, err := g.GetAsciiArmor()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return asciiArmor
}

func (g *GnuPGSignOptions) SetAsciiArmor(asciiArmor bool) (err error) {
	g.AsciiArmor = asciiArmor

	return nil
}

func (g *GnuPGSignOptions) SetDetachedSign(detachedSign bool) (err error) {
	g.DetachedSign = detachedSign

	return nil
}
