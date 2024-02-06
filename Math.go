package asciichgolangpublic

import "math"

type MathService struct{}

func Math() (m *MathService) {
	return NewMathService()
}

func NewMathService() (m *MathService) {
	return new(MathService)
}

func (m *MathService) MaxInt(integers ...int) (maxValue int) {
	maxValue = math.MinInt

	for _, i := range integers {
		if maxValue < i {
			maxValue = i
		}
	}

	return maxValue
}
