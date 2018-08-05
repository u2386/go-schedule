package schedule

import "reflect"

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
	in := make([]reflect.Value, len(args))
	for i := range args {
		in[i] = reflect.ValueOf(args[i])
	}
	s.addJob(&job{
		f:        fv,
		args:     in,
		interval: 1,
	})
}
