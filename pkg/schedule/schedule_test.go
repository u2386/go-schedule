package schedule

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

var (
	out       = make(chan interface{}, 1)
	scheduler Scheduler

	tests = []struct {
		fn     interface{}
		args   []interface{}
		expect interface{}
	}{
		{func(s string) {
			out <- fmt.Sprint(s)
		}, []interface{}{"WithArgs"}, "WithArgs"},
		{func() {
			out <- "WithoutArgs"
		}, nil, "WithoutArgs"},
	}
)

func TestJob(t *testing.T) {
	j := &job{
		f: reflect.ValueOf(func(s string) {
			out <- fmt.Sprint(s)
		}),
		args: []reflect.Value{reflect.ValueOf("Test")},
	}
	j.Run()
	assert.Equal(t, <-out, "Test", "they should be equal")
}

func TestSchedule(t *testing.T) {
	for _, test := range tests {
		scheduler.Do(test.fn, test.args...)
		scheduler.RunOnce()
		assert.Equal(t, <-out, test.expect, "they should be equal")
	}
}
