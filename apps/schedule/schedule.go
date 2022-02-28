package schedule

import (
	"fmt"
	"github.com/getevo/evo-ng"
	"github.com/getevo/evo-ng/lib/hash"
	"sync"
	"time"
)

type Job struct {
	Name         string         `json:"name"`
	Fn           func() error   `json:"-"`
	OnError      func(error)    `json:"-"`
	Blocking     bool           `json:"blocking"`
	NextInvoke   time.Time      `json:"next_invoke"`
	LastInvoke   time.Time      `json:"last_invoke"`
	Recurring    *time.Duration `json:"recurring"`
	ExecTime     int64          `json:"exec_time"`
	AvgExecTime  int64          `json:"avg_exec_time"`
	SingleNode   bool           `json:"single_node"`
	LastError    error          `json:"last_error"`
	Paused       bool           `json:"paused"`
	Running      bool           `json:"running"`
	Instances    int            `json:"instances"`
	MaxInstances int            `json:"max_instances"`
	Limiter      *evo.Limiter   `json:"limiter"`
}

var mu sync.Mutex
var jobs = map[string]*Job{}
var registered = false
var ticker = time.Duration(1 * time.Second)

func Jobs() map[string]*Job {
	return jobs
}

func SetPrecision(duration time.Duration) {
	ticker = duration
}

func Register() error {
	if registered {
		return nil
	}
	registered = true
	go func() {
		for {
			var now = time.Now()
			for idx, job := range jobs {
				if job.Paused || (job.Blocking && job.Running) || (job.MaxInstances > 0 && job.MaxInstances < job.Instances) {
					continue
				}
				if job.Limiter != nil && !(*job.Limiter).TryAcquire() {
					return
				}
				if job.NextInvoke.Before(now) {
					job := jobs[idx]
					go func() {
						mu.Lock()
						job.LastInvoke = now
						job.Instances = job.Instances + 1
						job.Running = true
						mu.Unlock()
						var err = job.Fn()
						mu.Lock()
						job.LastError = err
						if job.OnError != nil {
							job.OnError(err)
						}
						job.Instances = job.Instances - 1
						if job.Blocking {
							job.Running = false
						} else {
							if job.Instances < 1 {
								job.Running = false
							}
						}
						job.ExecTime = time.Now().UnixMilli() - job.LastInvoke.UnixMilli()
						if job.AvgExecTime > 0 {
							job.AvgExecTime = (job.ExecTime + job.AvgExecTime) / 2
						} else {
							job.AvgExecTime = job.ExecTime
						}
						if job.Recurring != nil {
							job.NextInvoke = now.Add(*job.Recurring)
							mu.Unlock()
						} else {
							mu.Unlock()
							job.Stop()
						}
					}()
				}
			}
			time.Sleep(ticker)
		}
	}()
	return nil
}

func NewSchedule(name string, fn func() error) *Job {
	return &Job{
		Name: name,
		Fn:   fn,
	}
}

func (j *Job) SetOnError(fn func(error)) *Job {
	j.OnError = fn
	return j
}

func (j *Job) SetRecurring(every time.Duration) *Job {
	j.Recurring = &every
	j.NextInvoke = time.Now().Add(every)
	return j
}

func (j *Job) SetNextInvoke(next interface{}) *Job {
	switch v := next.(type) {
	case string:
		var duration, _ = time.ParseDuration(v)
		j.NextInvoke = time.Now().Add(duration)
	case time.Duration:
		j.NextInvoke = time.Now().Add(v)
	case time.Time:
		j.NextInvoke = v
	}

	return j
}

func (j *Job) SetMaxInstances(instances int) *Job {
	j.MaxInstances = instances
	return j
}

func (j *Job) SetBlocking(blocking bool) *Job {
	j.Blocking = blocking
	return j
}

func (j *Job) SetLimiter(limiter evo.Limiter) *Job {
	j.Limiter = &limiter
	return j
}

func (j *Job) Pause() {
	j.Paused = true
}

func (j *Job) Start() error {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := jobs[j.Name]; ok {
		if jobs[j.Name].Paused {
			jobs[j.Name].Paused = false
			return nil
		}
		return fmt.Errorf("another job with same name exists")
	}
	jobs[j.Name] = j
	return nil
}

func (j *Job) Stop() {
	mu.Lock()
	defer mu.Unlock()
	delete(jobs, j.Name)
}

func SetTimeout(fn func() error, duration time.Duration) *Job {
	var job = NewSchedule(hash.UUID(), fn).SetNextInvoke(time.Now().Add(duration)).SetBlocking(true)
	job.Start()
	return job
}

func SetInterval(fn func() error, duration time.Duration) *Job {
	var job = NewSchedule(hash.UUID(), fn).SetRecurring(duration).SetBlocking(true)
	job.Start()
	return job
}
