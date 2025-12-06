package util

import (
	"time"
)

func TimeoutBySeconds(seconds time.Duration) {
	time.Sleep(time.Second * seconds)
}
