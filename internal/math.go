package internal

import (
	"math"
	"slices"

	"github.com/samber/lo"
	"golang.org/x/exp/constraints"
)

type Summable interface {
	constraints.Float | constraints.Integer
}

func DotProduct[T Summable](v1 []T, v2 []T) (T, error) {
	if len(v1) != len(v2) {
		return 0, ErrVectorLenMismtach
	}
	return lo.Reduce(v1, func(agg T, value T, i int) T {
		return value*v2[i] + agg
	}, 0), nil
}

func Magnitude[T Summable](ts []T) float64 {
	sum := lo.SumBy(ts, func(t T) T {
		return t * t
	})
	return math.Sqrt(float64(sum))
}

func CosineDistance[T Summable](v1 []T, v2 []T) (float64, error) {
	if slices.Equal(v1, v2) {
		return 0, nil
	}
	dp, err := DotProduct(v1, v2)
	if err != nil {
		return 0, err
	}
	m1 := Magnitude(v1)
	m2 := Magnitude(v2)
	if m1 == 0 || m2 == 0 {
		return 0, nil
	}
	return 1 - (float64(dp) / (m1 * m2)), nil
}
