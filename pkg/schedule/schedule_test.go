package schedule

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

var (
	out      = make(chan interface{}, 1)
	scheduler Scheduler
)

func runWithReturn(s string, o chan<- interface{}) {
	o <- fmt.Sprint(s)
}

func TestJob(t *testing.T) {
	args := []interface{}{"Test", out}
	in := make([]reflect.Value, len(args))
	for i := range args {
		in[i] = reflect.ValueOf(args[i])
	}
	j := &job{
		f:    reflect.ValueOf(runWithReturn),
		args: in,
	}
	j.Run()
	assert.Equal(t, <-out, "Test", "they should be equal")
}

func TestSchedule(t *testing.T) {
	scheduler.Do(runWithReturn, "Test", out)
	scheduler.RunOnce()
	assert.Equal(t, <-out, "Test", "they should be equal")
}
