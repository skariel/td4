// Package utils just some general utils for everybody
package utils

import (
	"log"
	"os"
	"time"
)

// DoEvery execute the given function in parameters every given duration. Also executes once at start
func DoEvery(d time.Duration, fn func()) {
	fn()

	for range time.Tick(d) {
		fn()
	}
}

// LoggedGetEnv logs var name, var val and returns it
func LoggedGetEnv(varName string) string {
	varVal := os.Getenv("TD4_CERTIFICATE_FILE_PATH")
	log.Printf("env %v = %v", varName, varVal)

	return varVal
}
