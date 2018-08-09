package schedule

import (
	"reflect"
	"time"
)

type job struct {
	f         reflect.Value
	args      []reflect.Value
	interval  uint
	unit      time.Duration
	scheduler *Scheduler
	nextTime  time.Time
	lastTime  time.Time
}

func (j *job) run() {
	j.lastTime = time.Now()
	j.f.Call(j.args)
}

// Every schedules a new periodic job.
func (j *job) Every(interval uint) *job {
	j.interval = interval
	return j
}

func (j *job) Seconds() {
	j.unit = time.Second
}

func (j *job) Minutes() {
	j.unit = time.Minute
}

func (j *job) Hours() {
	j.unit = time.Hour
}
