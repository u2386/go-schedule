package schedule

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
		go j.run()
	}
}

// RunAll is like RunOnce but runs forever.
func (s *Scheduler) RunAll() {
	for {
		for _, j := range s.jobs {
			go j.run()
		}
	}
}

// Every schedules a new periodic job.
func (s *Scheduler) Every(interval int) *job {
	if interval <= 0 {
		panic("interval should greater than 0.")
	}
	return &job{
		interval:  interval,
		scheduler: s,
	}
}
