package utils

import "time"

// TimeIsZero actually returns true if time is zero
func TimeIsZero(t time.Time) bool {
	return t.UnixNano() == 0 || t.IsZero()
}
