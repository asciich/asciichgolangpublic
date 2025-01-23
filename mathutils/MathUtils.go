package mathutils

import "math"

func MaxInt(integers ...int) (maxValue int) {
	maxValue = math.MinInt

	for _, i := range integers {
		if maxValue < i {
			maxValue = i
		}
	}

	return maxValue
}
