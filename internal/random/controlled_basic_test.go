package random

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBasicControl_Execute(t *testing.T) {
	UseTestSeed()

	ctrl := &BasicControl{
		ID:     "hello_world",
		Chance: 1.0,
	}

	assert.True(t, ctrl.Execute())
	ctrl.Chance = 0
	assert.False(t, ctrl.Execute())
}

func TestBasicControl_ToAnalytics(t *testing.T) {
	var (
		exID   = "hello_world"
		chance = 0.42069
	)

	ctrl := &BasicControl{
		ID:     exID,
		Chance: chance,
	}

	assert.Equal(t, exID, ctrl.ToAnalytics().ID)
	assert.Equal(t, chance, ctrl.ToAnalytics().Chance)
	assert.Equal(t, 0, ctrl.ToAnalytics().CallCount)

	ctrl.Execute()
	assert.Equal(t, 1, ctrl.ToAnalytics().CallCount)

	anl := ctrl.ToAnalytics()
	assert.True(t, anl.TrueCount == 1 || anl.FalseCount == 1)
}
