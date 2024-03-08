package random

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"runtime/debug"
	"strings"
	"time"
)

var (
	globalRandomSeed *rand.Rand
	seedFromTime     = false
	seedTime         time.Time
)

var keyTable = [6]string{
	"Y",
	"M",
	"d",
	"h",
	"m",
	"s",
}

var valTable = [6]int64{
	86400 * 365,
	86400 * 30,
	86400,
	3600,
	60,
	1,
}

func using() *rand.Rand {
	if globalRandomSeed == nil {
		UseTestSeed()
	}

	if seedFromTime {
		now := time.Now()
		if now.Unix()-seedTime.Unix() > 3600 {
			UseTimeBasedSeed()
		}
	}

	return globalRandomSeed
}

// UseTestSeed provides a stable basis for testing with a seed of 0
func UseTestSeed() {
	seedFromTime = false
	seedTime = time.Now()
	UseSeed(0)
}

// UseSeed declares what seed our RNG system should use
func UseSeed(seed int64) {
	globalRandomSeed = rand.New(rand.NewSource(seed))
}

// UseTimeBasedSeed grabs a seed from `time.Now().UnixNano()`
func UseTimeBasedSeed() {
	seedFromTime = true
	seedTime = time.Now()

	UseSeed(seedTime.UnixNano())
}

// Float returns a random float64
func Float() float64 {
	return using().Float64()
}

// Float32 returns a random float32
func Float32() float32 {
	return using().Float32()
}

// Float64 returns a random float64
func Float64() float64 {
	return using().Float64()
}

// Float64Range takes a range between min & max and adds a Float64() to it
func Float64Range(min, max float64) float64 {
	if min >= max {
		panic("min should not be bigger than max")
	}

	return min + rand.Float64()*(max-min)
}

// Int64 returns a random int64 from 0...n
func Int64(n int64) int64 {
	// just in case bby
	if n <= 0 {
		log.Println("Using Int64() with a <=0 number!")
		debug.PrintStack()
		return 1
	}

	return using().Int63n(n)
}

// Int64Range returns a number between min/max
func Int64Range(min, max int64) int64 {
	if min >= max {
		panic("min should not be bigger than max")
	}

	return Int64(max-min) + min
}

// UInt32 returns a number between 0 and n
func UInt32(n uint32) uint32 {
	ret := using().Int63n(int64(n))

	return uint32(ret % math.MaxUint32)
}

// Int32 returns a number between 0 and n
func Int32(n int32) int32 {
	return using().Int31n(n)
}

// Int32Range returns a number between 0 and n
func Int32Range(min, max int32) int32 {
	if min >= max {
		panic("min should not be bigger than max")
	}

	return Int32(max-min) + min
}

// Int returns a number between 0 and n
func Int(n int) int {
	return int(Int64(int64(n)))
}

// PickString returns one of the choices provided (and trims the string)
func PickString(choices []string) string {
	choiceLen := len(choices)

	if choiceLen <= 0 {
		return ""
	}

	return strings.TrimSpace(choices[Int(choiceLen)])
}

// RandomlyTrue decides what is true weighted by chance.
// Chance is between [0,1] where 0.0 is false always, 1.0 is true always
func RandomlyTrue(chance float64) bool {
	if chance >= 1.0 {
		return true
	}

	if chance <= 0.0 {
		return false
	}

	return chance >= using().Float64()
}

// legacyRandomlyTrue decides what is true, weighted to a specific truth
// must be between 0.0 and 1.0 (non-inclusive), so 0.0 will be
// false always, 1.0 will be true always.
func legacyRandomlyTrue(weightOfTruth float64) bool {
	// if true is 1.0 or above it's 100% the case.
	if weightOfTruth >= 1.0 {
		return true
	}

	// if true is 0.0 or below it's never going to happen.
	if weightOfTruth <= 0.0 {
		return false
	}

	// TODO simplify while maintaining the cohesion of this algorithm thing

	totalWeight := 1.0
	randomWeight := using().Float64() * totalWeight
	accruedWeight := 0.0

	choices := map[bool]float64{
		true:  weightOfTruth,
		false: 1.0 - weightOfTruth,
	}

	for val, chance := range choices {
		accruedWeight += chance

		if accruedWeight >= randomWeight {
			return val
		}
	}

	log.Printf("WARN: Shouldn't get here!\n")
	return false
}

// MakeTime generates a random "duration" between the
// min and max seconds provided and returns a string
// of the conversion.
// (so: 3600,3601 would return 1h,1h 1s respectively)
func MakeTime(min, max int64) string {

	move := min + using().Int63n(max-min)
	output := ""

	for pos, val := range valTable {
		key := keyTable[pos]

		if move/val > 0 {
			output += fmt.Sprintf(" %d%s", move/val, key)
			move -= val * (move / val)
		}
	}

	return strings.TrimSpace(output)
}
