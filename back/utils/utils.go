// Package utils just some general utils for everybody
package utils

import "time"

// DoEvery execute the given function in parameters every given duration. Also executes once at start
func DoEvery(d time.Duration, fn func()) {
	fn()

	for range time.Tick(d) {
		fn()
	}
}
