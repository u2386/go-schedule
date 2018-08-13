package schedule

import (
	"reflect"
	"time"
)

// States of Job
const (
	Waiting uint = 0
	Running
)

type job struct {
	f         reflect.Value
	args      []reflect.Value
	interval  uint
	unit      time.Duration
	scheduler *Scheduler
	nextTime  time.Time
	lastTime  time.Time
	state     uint
}

func (j *job) run() {
	j.lastTime = time.Now()
	j.state = Running
	j.f.Call(j.args)
	j.state = Waiting
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

func (j *job) canRun() bool {
	return j.state == Waiting && time.Now().After(j.nextTime)
}
