package schedule

import (
	"fmt"
	"reflect"
	"runtime"
	"time"
)

// Scheduler is a realtime job-arrangement manager.
type Scheduler struct {
	jobs []*job
}

// RunOnce runs all jobs.
func (s *Scheduler) RunOnce() {
	for _, j := range s.jobs {
		if j.state == Waiting {
			s.schedule(j)
		}
		if j.canRun() {
			go j.run()
		}
	}
}

func (s *Scheduler) schedule(j *job) {
	if j.nextTime.After(time.Now()) {
		return
	}
	j.nextTime = time.Now().Add(time.Duration(j.interval) * j.unit)
}

// Do accepts function pointer and it's arguments for scheduler to run.
func (s *Scheduler) Do(fn interface{}, args ...interface{}) (j *job) {
	fv := reflect.ValueOf(fn)
	if fv.Kind() != reflect.Func {
		panic("Argument fn must be a function type")
	}
	ft := fv.Type()
	in := make([]reflect.Value, ft.NumIn())
	if ft.NumIn() == 0 {
		in = nil
	} else {
		for i := range in {
			in[i] = reflect.ValueOf(args[i])
			if in[i].Kind() != ft.In(i).Kind() {
				fname := runtime.FuncForPC(fv.Pointer()).Name()
				panic(fmt.Sprintf("cannot use type %T as type %s in argument to %s", args[i], ft.In(i).Name(), fname))
			}
		}
	}
	j = &job{
		f:         fv,
		args:      in,
		scheduler: s,
		state:     Waiting,
	}
	s.jobs = append(s.jobs, j)
	return
}
