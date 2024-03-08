package utils

import (
	"log"
	"os"
)

// MustSetEnv is a wrapper around os.Setenv to set an environmental variable
// or panic if it fails.
func MustSetEnv(name, value string) {
	err := os.Setenv(name, value)
	if err != nil {
		log.Printf("Failed to SetEnv: %v\n", err)
		panic(err)
	}
}

// WithCICDEnvVars allows us to wrap around env variables definitely going
//
func WithCICDEnvVars(wrapper func()) {
	origCI := os.Getenv(EnvCI)
	origCD := os.Getenv(EnvCD)

	MustSetEnv(EnvCI, "1")
	MustSetEnv(EnvCD, "1")

	wrapper()

	MustSetEnv(EnvCI, origCI)
	MustSetEnv(EnvCD, origCD)
}
