package schedule

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var out = make(chan interface{}, 1)

func TestJob(t *testing.T) {
	j := &job{
		f: reflect.ValueOf(func(s string) {
			out <- fmt.Sprint(s)
		}),
		args: []reflect.Value{reflect.ValueOf("Test")},
	}
	j.run()
	assert.Equal(t, <-out, "Test", "they should be equal")
}

func TestSchedulePanic(t *testing.T) {
	var scheduler Scheduler
	defer func() {
		r := recover()
		if r == nil || !strings.HasPrefix(r.(string), "cannot use type int as type string in argument") {
			t.Errorf("should panic for argument type mismatch")
		}
	}()
	scheduler.Do(func(s string) {
		out <- fmt.Sprint(s)
	}, 1).Every(1)
}

var tests = []struct {
	fn     interface{}
	args   []interface{}
	expect interface{}
}{
	{
		func(s string) {
			out <- fmt.Sprint(s)
		}, []interface{}{"WithArgs"}, "WithArgs",
	},
	{
		func() {
			out <- "WithoutArgs"
		}, nil, "WithoutArgs",
	},
}

func TestSchedule(t *testing.T) {
	var scheduler Scheduler
	for _, test := range tests {
		scheduler.Do(test.fn, test.args...).Every(1)
		scheduler.RunOnce()
		assert.Equal(t, test.expect, <-out, "they should be equal")
	}
}

func TestScheduleUnits(t *testing.T) {
	var j *job
	var scheduler Scheduler

	scheduler.Do(func() {}).Every(1).Seconds()
	j = scheduler.jobs[len(scheduler.jobs)-1]
	assert.Equal(t, time.Second, j.unit, "they should be equal")

	scheduler.Do(func() {}).Every(1).Minutes()
	j = scheduler.jobs[len(scheduler.jobs)-1]
	assert.Equal(t, time.Minute, j.unit, "they should be equal")

	scheduler.Do(func() {}).Every(1).Hours()
	j = scheduler.jobs[len(scheduler.jobs)-1]
	assert.Equal(t, time.Hour, j.unit, "they should be equal")
}

func (s *Scheduler) mockSchedule(j *job, now time.Time) {
	j.nextTime = now.Add(time.Duration(j.interval) * j.unit)
}

func TestScheduleNextTime(t *testing.T) {
	var j *job
	var scheduler Scheduler
	now := time.Date(2018, 8, 9, 22, 40, 59, 0, time.UTC)

	j = &job{interval: 2, unit: time.Second}
	scheduler.mockSchedule(j, now)
	assert.Equal(t, time.Date(2018, 8, 9, 22, 41, 1, 0, time.UTC), j.nextTime)

	j = &job{interval: 1, unit: time.Minute}
	scheduler.mockSchedule(j, now)
	assert.Equal(t, time.Date(2018, 8, 9, 22, 41, 59, 0, time.UTC), j.nextTime)

	j = &job{interval: 1, unit: time.Hour}
	scheduler.mockSchedule(j, now)
	assert.Equal(t, time.Date(2018, 8, 9, 23, 40, 59, 0, time.UTC), j.nextTime)
}
