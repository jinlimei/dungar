package triggers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTriggerCallback_Matches(t *testing.T) {
	nonRgx := triggerCallback{
		Trigger: "hello",
		IsRegex: true,
	}

	assert.True(t, nonRgx.matches("hello"))
	assert.False(t, nonRgx.matches("world"))

	withRgx := triggerCallback{
		Trigger: "\\Wworld",
		IsRegex: true,
	}

	assert.True(t, withRgx.matches("hello world"))
	assert.True(t, withRgx.matches("banana worldstar"))
	assert.False(t, withRgx.matches("worldstar"))
	assert.False(t, withRgx.matches("world"))
}

func TestTriggerChoice_Matches(t *testing.T) {
	simple := triggerChoice{
		Trigger:  "mercy",
		IsRegex:  false,
		Response: "mercy main btw",
		cache:    nil,
	}

	assert.True(t, simple.matches("i just want to play mercy"))
}

func TestMultiTriggerCallback_Matches(t *testing.T) {
	// our ~ non-regex zone ~
	nonRegex := multiTriggerCallback{
		Triggers: []string{"woo", "world"},
		IsRegex:  false,
	}

	assert.True(t, nonRegex.matches("hello world"))
	assert.False(t, nonRegex.matches("swoon"))

	// our ~ regex zone ~
	withRegex := multiTriggerCallback{
		Triggers: []string{"\\Wworld"},
		IsRegex:  true,
	}

	assert.True(t, withRegex.matches("hello world"))
	assert.True(t, withRegex.matches("hello worldstar"))
	assert.False(t, withRegex.matches("world"))
}
