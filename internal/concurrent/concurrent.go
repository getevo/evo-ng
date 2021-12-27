package concurrent

import (
	context "context"
	"time"
)

func Limit(args ...interface{}) error {
	var ctx = context.Background()
	var timeout = 0 * time.Hour
	for _, arg := range args {
		switch arg.(type) {
		case context.Context:
			ctx = arg.(context.Context)
		case time.Duration:
			timeout = arg.(time.Duration)
		}
	}

	return nil
}

func Retry(args ...interface{}) error {
	var ctx = context.Background()
	var timeout = 0 * time.Hour
	for _, arg := range args {
		switch arg.(type) {
		case context.Context:
			ctx = arg.(context.Context)
		case time.Duration:
			timeout = arg.(time.Duration)
		}
	}

	return nil
}

func First(args ...interface{}) error {
	var ctx = context.Background()
	var timeout = 0 * time.Hour
	for _, arg := range args {
		switch arg.(type) {
		case context.Context:
			ctx = arg.(context.Context)
		case time.Duration:
			timeout = arg.(time.Duration)
		}
	}

	return nil
}

func All(args ...interface{}) error {
	var ctx = context.Background()
	var timeout = 0 * time.Hour
	for _, arg := range args {
		switch arg.(type) {
		case context.Context:
			ctx = arg.(context.Context)
		case time.Duration:
			timeout = arg.(time.Duration)
		}
	}

	return nil
}
