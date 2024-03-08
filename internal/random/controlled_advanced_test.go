package random

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
	"time"
)

func TestAdvancedControl_Execute(t *testing.T) {
	var (
		ctrl AdvancedControl
	)

	for k := 0; k < 100; k++ {
		UseTestSeed()

		ctrl = AdvancedControl{
			ID:            "hello_world",
			DefaultChance: 1.0,
		}

		ctrl.lastTimeTrue = time.Time{}
		ctrl.lastTimeCalled = time.Time{}

		assert.Equal(t, 0.0, ctrl.currentChance,
			"Default current chance should be 0.0")
		assert.True(t, ctrl.Execute(),
			"Control should have been true when executing (init is auto-called)")
		assert.NotEqual(t, 0, ctrl.lastTimeTrue.Second(),
			"last time true should be at the start of the epoch")

		ctrl.Init()

		assert.Equal(t, ctrl.DefaultChance, ctrl.currentChance,
			"Default & current Chances should be the same")
		assert.True(t, ctrl.Execute(),
			"Control should have been true when executing.")
		assert.NotEqual(t, 0, ctrl.lastTimeTrue.Second(),
			"last time true should be now(ish)")

		ctrl.lastTimeTrue = time.Now().Add(-13 * time.Minute)
		ctrl.MaximumOccurrence = 14 * time.Minute

		assert.False(t, ctrl.Execute(),
			"Control was executed before MaximumTime so should be false")

		ctrl.lastTimeTrue = time.Now().Add(-15 * time.Minute)
		assert.True(t, ctrl.Execute(),
			"Control was executed after MaximumTime so should be true")

		// reset to be always false so let's get the game running again
		ctrl.DefaultChance = 0.0
		ctrl.currentChance = 0.0

		ctrl.lastTimeTrue = time.Now().Add(-16 * time.Minute)
		ctrl.MaximumOccurrence = 0
		ctrl.MinimumOccurrence = time.Duration(15 * time.Minute)

		assert.True(t, ctrl.Execute(),
			"Control should execute as true since we're passed the MinimumOccurrence")
	}
}

func TestAdvancedControl_ExecuteIncrement(t *testing.T) {
	now := time.Now()

	for k := 0; k < 5; k++ {
		UseTestSeed()

		ctrl := AdvancedControl{
			ID:                "hello_world",
			IncrementIncrease: 0.25,
			IncreaseDuration:  1 * time.Minute,
			DefaultChance:     0.00,
		}

		ctrl.lastTimeTrue = time.Now()
		ctrl.lastTimeCalled = time.Now()

		ctrl.Execute()

		assert.Equal(t, ctrl.DefaultChance, ctrl.currentChance)
		assert.Equal(t, 0.0, ctrl.currentChance)

		ctrl.lastTimeTrue = now.Add(-1 * time.Minute)
		ctrl.lastTimeCalled = now.Add(-1 * time.Minute)

		ctrl.Execute()

		assert.NotEqual(t, ctrl.DefaultChance, ctrl.currentChance)
		assert.NotEqual(t, 0.0, ctrl.currentChance)
		assert.Truef(t, closeEnough(0.25, ctrl.currentChance),
			"expected (%f) and current (%f) were not close enough", 0.25, ctrl.currentChance)

		ctrl.lastTimeCalled = now.Add(-2 * time.Minute)
		ctrl.lastTimeTrue = now.Add(-2 * time.Minute)
		ctrl.currentChance = ctrl.DefaultChance

		ctrl.Execute()

		assert.NotEqual(t, ctrl.DefaultChance, ctrl.currentChance)
		assert.NotEqual(t, 0.0, ctrl.currentChance)
		assert.True(t, closeEnough(0.5, ctrl.currentChance))

		ctrl.lastTimeCalled = ctrl.lastTimeCalled.Add(-1 * time.Minute)
		ctrl.lastTimeTrue = ctrl.lastTimeTrue.Add(-1 * time.Minute)

		res := ctrl.Execute()

		// We _can_ return true here on chance so, let's address that scenario.
		if res {
			assert.Equal(t, ctrl.DefaultChance, ctrl.currentChance)
			assert.Equal(t, 0.0, ctrl.currentChance)
		} else {
			assert.NotEqualf(t, ctrl.DefaultChance, ctrl.currentChance,
				"did not expect default(%f) but got current(%f)", ctrl.DefaultChance, ctrl.currentChance)
			assert.NotEqualf(t, 0.0, ctrl.currentChance,
				"current chance(%f) should not be zero", ctrl.currentChance)
			assert.True(t, closeEnough(0.75, ctrl.currentChance))
		}
	}
}

func TestAdvancedControl_ToAnalytics(t *testing.T) {
	var (
		exID = "foo"
		chance = 0.42069
	)

	ctrl := &AdvancedControl{
		ID:            exID,
		DefaultChance: chance,
	}

	assert.False(t, ctrl.isReady)
	assert.NotEqual(t, chance, ctrl.currentChance)

	ctrl.Init()

	assert.True(t, ctrl.isReady)
	assert.Equal(t, chance, ctrl.currentChance)

	anl := ctrl.ToAnalytics()
	assert.Equal(t, chance, anl.Chance)
	assert.Equal(t, exID, anl.ID)
	assert.Equal(t, 0, anl.CallCount)
	assert.Equal(t, 0, anl.TrueCount)
	assert.Equal(t, 0, anl.FalseCount)

	ctrl.Execute()
	anl = ctrl.ToAnalytics()
	assert.Equal(t, chance, anl.Chance)
	assert.Equal(t, exID, anl.ID)
	assert.Equal(t, 1, anl.CallCount)
	assert.True(t, anl.TrueCount == 1 || anl.FalseCount == 1)
}

func TestCloseEnough(t *testing.T) {
	assert.False(t, closeEnough(0.25, 0.256))
	assert.False(t, closeEnough(0.25, 0.2506))
	assert.False(t, closeEnough(0.25, 0.25006))

	assert.True(t, closeEnough(0.25, 0.250006))
	assert.True(t, closeEnough(0.25, 0.2500006))
	assert.True(t, closeEnough(0.25, 0.25000006))
	assert.True(t, closeEnough(0.25, 0.250000006))
	assert.True(t, closeEnough(0.25, 0.2500000006))
	assert.True(t, closeEnough(0.25, 0.25000000006))
	assert.True(t, closeEnough(0.25, 0.250000000006))
	assert.True(t, closeEnough(0.25, 0.2500000000006))
	assert.True(t, closeEnough(0.25, 0.25000000000006))
	assert.True(t, closeEnough(0.25, 0.250000000000006))
	assert.True(t, closeEnough(0.25, 0.2500000000000006))
}

func closeEnough(x, y float64) bool {
	var (
		x1 = float64(x * 10_000)
		y1 = float64(y * 10_000)
	)

	return math.Round(x1) == math.Round(y1)
}
