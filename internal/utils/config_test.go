package utils

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMustUseEnvVars(t *testing.T) {
	origCI := os.Getenv(EnvCI)
	origCD := os.Getenv(EnvCD)

	MustSetEnv(EnvCI, "")
	MustSetEnv(EnvCD, "")
	assert.False(t, MustUseEnvVars())

	MustSetEnv(EnvCI, "1")
	MustSetEnv(EnvCD, "1")
	assert.True(t, MustUseEnvVars())

	MustSetEnv(EnvCI, origCI)
	MustSetEnv(EnvCD, origCD)
}

