package util

import (
	"time"
)

var LastKnownStableTime time.Time

func GetUptime() time.Duration {
	if LastKnownStableTime.IsZero() {
		return 0
	}
	return time.Since(LastKnownStableTime)
}

func SetUptime() {
	LastKnownStableTime = time.Now()
}
