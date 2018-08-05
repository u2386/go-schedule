package schedule

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	out       = make(chan interface{}, 1)
	scheduler Scheduler
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
	for _, test := range tests {
		scheduler.Do(test.fn, test.args...)
		scheduler.RunOnce()
		assert.Equal(t, <-out, test.expect, "they should be equal")
	}
}

func TestSchedulePanic(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil || !strings.HasPrefix(r.(string), "cannot use type int as type string in argument") {
			t.Errorf("should panic for argument type mismatch")
		}
	}()
	scheduler.Do(func(s string) {
		out <- fmt.Sprint(s)
	}, 1)
	scheduler.RunOnce()
}
