package schedule

import "time"

type Job struct {
	Fn         func() error
	Blocking   bool
	NextInvoke time.Time
	LastInvoke time.Time
	Recurring  *time.Duration
	SingleNode bool
	//DistLock        dist.Lock
}

func Register() error {
	go func() {
		for {

			time.Sleep(1 * time.Second)
		}
	}()
	return nil
}

func NewSchedule() *Job {

}
