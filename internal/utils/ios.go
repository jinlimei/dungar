package utils

import (
	"log"
	"os"
	"path"
	"strings"
)

// GetExecutablePath returns the directory path from where we are executing.
func GetExecutablePath() string {
	exe, err := os.Executable()
	if err != nil {
		log.Fatalln(err)
	}

	exe = strings.ReplaceAll(exe, "\\", "/")
	return path.Dir(exe)
}

// GetCurrentDir returns the current working directory for where we're executing
func GetCurrentDir() string {
	dio, err := os.Getwd()

	if err != nil {
		log.Fatalln(err)
	}

	dio = strings.ReplaceAll(dio, "\\", "/")

	return dio
}
