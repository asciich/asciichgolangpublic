package randomgenerator

import (
	"math/rand"

	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func GetRandomString(lenght int) (string, error) {
	if lenght <= 0 {
		return "", tracederrors.TracedErrorf("Invalid lenght '%d' to generate random string", lenght)
	}

	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, lenght)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b), nil
}
