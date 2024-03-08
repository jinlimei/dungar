package random

import "time"

func timeIsZero(t time.Time) bool {
	return t.UnixNano() == 0 || t.IsZero()
}
