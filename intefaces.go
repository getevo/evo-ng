package evo

import "time"

type Limiter interface {
	TryAcquireDuration(duration time.Duration) bool
	TryAcquire() bool
}
