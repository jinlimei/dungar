package random

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func isAround(num int, around int, spread int) bool {
	return (around-spread) < num && (around+spread) > num
}

type randomlyTrueTestCase struct {
	weight   float64
	trueLim  int
	falseLim int
	spread   int
}

type rngTrue func(c float64) bool

func TestRandomlyTrue(t *testing.T) {
	UseTestSeed()

	assert.True(t, legacyRandomlyTrue(1.0))
	assert.False(t, legacyRandomlyTrue(0.0))
	assert.True(t, RandomlyTrue(1.0))
	assert.False(t, RandomlyTrue(0.0))

	var (
		whenTrue    = 0
		whenFalse   = 0
		truePasses  = false
		falsePasses = false
	)

	testWeight := func(fnc rngTrue, weight float64, trueLim, falseLim, spread int) bool {
		whenTrue = 0
		whenFalse = 0

		for k := 0; k < 1_000_000; k++ {
			if fnc(weight) {
				whenTrue++
			} else {
				whenFalse++
			}
		}

		truePasses = isAround(whenTrue, trueLim, spread)
		falsePasses = isAround(whenFalse, falseLim, spread)

		assert.True(t, truePasses,
			fmt.Sprintf("whenTrue failed (weight=%f, trueLim=%d, whenTrue=%d, spread=%d)",
				weight, trueLim, whenTrue, spread))
		assert.True(t, falsePasses,
			fmt.Sprintf("whenFalse failed (weight=%f, trueLim=%d, whenFalse=%d, spread=%d)",
				weight, falseLim, whenFalse, spread))

		log.Printf("legacyRandomlyTrue: weight=%f, whenTrue=%d, whenFalse=%d, spread=%d\n",
			weight, whenTrue, whenFalse, spread)

		return truePasses && falsePasses
	}

	testCases := []randomlyTrueTestCase{
		{0.10, 100_000, 900_000, 1_500},
		{0.01, 10_000, 990_000, 150},
		{0.001, 1000, 999_000, 75},
		{0.0001, 100, 999_900, 75},
		// More uncertainty around here
		{0.00001, 10, 999_990, 10},
	}

	funcs := map[string]rngTrue{
		"RandomlyTrue":       RandomlyTrue,
		"legacyRandomlyTrue": legacyRandomlyTrue,
	}

	for name, fnc := range funcs {
		log.Printf("RUNNING WITH '%s'\n", name)

		for _, testCase := range testCases {
			UseTestSeed()
			testWeight(fnc, testCase.weight, testCase.trueLim, testCase.falseLim, testCase.spread)
		}

		log.Printf("DONE WITH '%s'\n", name)
	}
}

func BenchmarkLegacyRandomlyTrue(b *testing.B) {
	UseTestSeed()

	for k := 0; k < b.N; k++ {
		legacyRandomlyTrue(0.50)
	}
}

func BenchmarkRandomlyTrue(b *testing.B) {
	UseTestSeed()

	for k := 0; k < b.N; k++ {
		RandomlyTrue(0.50)
	}
}

func TestRandomlyTrueAsACoinFlip(t *testing.T) {
	whenTrue := 0
	whenFalse := 0

	for i := 0; i < 50000; i++ {
		if legacyRandomlyTrue(0.50) {
			whenTrue++
		} else {
			whenFalse++
		}
	}

	log.Printf("legacyRandomlyTrue had %d TRUE and %d FALSE\n", whenTrue, whenFalse)

	assert.True(t, isAround(whenTrue, 25000, 1500), strconv.Itoa(whenTrue))
	assert.True(t, isAround(whenTrue, 25000, 1500), strconv.Itoa(whenFalse))
}

func TestInt64Range(t *testing.T) {
	UseTestSeed()

	val := Int64Range(0, 100)
	assert.True(t, val < 100)
	assert.True(t, val > 0)

	for i := 0; i < 100; i++ {
		val = Int64Range(6000, 25000)
		assert.True(t, val > 6000)
		assert.True(t, val < 25000)
	}
}

func TestFloat64Range(t *testing.T) {
	UseTestSeed()

	assert.NotPanics(t, func() {
		Float64Range(0, 0.3)
		Float64Range(-0, 0.3)
		Float64Range(-0.3, 0)
	})
}

func TestMakeTime(t *testing.T) {
	UseTestSeed()

	assert.Equal(t, "3s", MakeTime(1, 10))
	assert.Equal(t, "1d 3h 29m 12s", MakeTime(89900, 99600))
}
