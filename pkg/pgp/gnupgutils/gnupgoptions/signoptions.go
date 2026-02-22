package gnupgoptions

import "github.com/asciich/asciichgolangpublic/pkg/logging"

type SignOption struct {
	DetachedSign bool
	AsciiArmor   bool
}

func NewGnuPGSignOptions() (g *SignOption) {
	return new(SignOption)
}

func (g *SignOption) GetAsciiArmor() (asciiArmor bool, err error) {

	return g.AsciiArmor, nil
}

func (g *SignOption) GetDetachedSign() (detachedSign bool, err error) {

	return g.DetachedSign, nil
}

func (g *SignOption) MustGetAsciiArmor() (asciiArmor bool) {
	asciiArmor, err := g.GetAsciiArmor()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return asciiArmor
}

func (g *SignOption) SetAsciiArmor(asciiArmor bool) (err error) {
	g.AsciiArmor = asciiArmor

	return nil
}

func (g *SignOption) SetDetachedSign(detachedSign bool) (err error) {
	g.DetachedSign = detachedSign

	return nil
}
