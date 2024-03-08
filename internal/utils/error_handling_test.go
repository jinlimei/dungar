package utils

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNonHaltingError(t *testing.T) {
	NonHaltingError("", nil)
	NonHaltingError("", errors.New("test"))
}

func TestHaltingError(t *testing.T) {
	assert.Panics(t, func() {
		HaltingError("", errors.New("hihi"))
	})

	assert.NotPanics(t, func() {
		HaltingError("", nil)
	})
}
