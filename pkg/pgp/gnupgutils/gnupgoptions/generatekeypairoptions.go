package gnupgoptions

import (
	"slices"

	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type GenerateKeyPairOptions struct {
	Name    string
	Comment string
	Email   string

	// Set RSA bit size
	RSABits int
}

func (g *GenerateKeyPairOptions) GetName() (string, error) {
	if g.Name == "" {
		return "", tracederrors.TracedError("Name not set")
	}

	return g.Name, nil
}

func (g *GenerateKeyPairOptions) GetEmail() (string, error) {
	if g.Email == "" {
		return "", tracederrors.TracedError("Email not set")
	}

	return g.Email, nil
}

func (g *GenerateKeyPairOptions) GetComment() (string, error) {
	if g.Comment == "" {
		return "", tracederrors.TracedError("Comment not set")
	}

	return g.Comment, nil
}

func (g *GenerateKeyPairOptions) GetRSABits() (int, error) {
	if g.RSABits == 0 {
		return 0, tracederrors.TracedError("RSABits not set")
	}

	rsaBits := g.RSABits

	allowedValues := []int{1024, 2048, 4096}
	if !slices.Contains(allowedValues, rsaBits) {
		return 0, tracederrors.TracedErrorf("Invalid RSABits: %d. Allowed values: %v", rsaBits, allowedValues)
	}

	return rsaBits, nil
}