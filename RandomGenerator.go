package asciichgolangpublic

import (
	"math/rand"

	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type RandomGeneratorService struct{}

func NewRandomGeneratorService() (randomGenerator *RandomGeneratorService) {
	return new(RandomGeneratorService)
}

func RandomGenerator() (randomGenerator *RandomGeneratorService) {
	return NewRandomGeneratorService()
}

func (r *RandomGeneratorService) GetRandomString(lenght int) (random string, err error) {
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

func (r *RandomGeneratorService) MustGetRandomString(length int) (random string) {
	random, err := r.GetRandomString(length)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return random
}
