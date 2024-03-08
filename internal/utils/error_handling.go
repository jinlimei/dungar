package utils

import (
	"github.com/davecgh/go-spew/spew"
	"log"
	"runtime/debug"
)

// NonHaltingError is to record errors that don't require process halting, but
// still should be recorded and captured
func NonHaltingError(loc string, err error) {
	if err == nil {
		return
	}

	log.Print(spew.Sdump(err))
	// TODO sentry-like error notification
	//SentryError(err, map[string]string{"location": loc})

	debug.PrintStack()
	log.Printf("Error @ '%s': %v\n", loc, err.Error())
}

// HaltingError records errors that should be recorded & halt the process.
func HaltingError(loc string, err error) {
	if err == nil {
		return
	}

	spew.Dump(err)

	// TODO sentry-like error notification
	//SentryError(err, map[string]string{"location": loc})
	debug.PrintStack()
	log.Panicf("Error @ '%s' %v\n", loc, err)
}
