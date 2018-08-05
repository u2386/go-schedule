package schedule

import (
	"fmt"
	"reflect"
	"runtime"
)

// Scheduler is a realtime job-arrangement manager.
type Scheduler struct {
	jobs []*job
}

func (s *Scheduler) addJob(j *job) {
	s.jobs = append(s.jobs, j)
}

// RunOnce runs all jobs.
func (s *Scheduler) RunOnce() {
	for _, j := range s.jobs {
		go j.Run()
	}
}

// RunAll is like RunOnce but runs forever.
func (s *Scheduler) RunAll() {
	for {
		for _, j := range s.jobs {
			go j.Run()
		}
	}
}

// Do accepts function pointer and it's arguments for scheduler to run.
func (s *Scheduler) Do(fn interface{}, args ...interface{}) {
	fv := reflect.ValueOf(fn)
	if fv.Kind() != reflect.Func {
		panic("Argument fn must be a function type")
	}
	ft := fv.Type()
	if ft.NumIn() == 0 {
		s.addJob(&job{
			f:        fv,
			args:     nil,
			interval: 1,
		})
		return
	}
	in := make([]reflect.Value, ft.NumIn())
	for i := range in {
		in[i] = reflect.ValueOf(args[i])
		if in[i].Kind() != ft.In(i).Kind() {
			fname := runtime.FuncForPC(fv.Pointer()).Name()
			panic(fmt.Sprintf("cannot use type %T as type %s in argument to %s", args[i], ft.In(i).Name(), fname))
		}
	}
	s.addJob(&job{
		f:        fv,
		args:     in,
		interval: 1,
	})
}
