package schedule

import (
	"fmt"
	"reflect"
	"runtime"
)

type job struct {
	f         reflect.Value
	args      []reflect.Value
	interval  int
	unit      int
	scheduler *Scheduler
}

func (j *job) run() {
	j.f.Call(j.args)
}

// Do accepts function pointer and it's arguments for scheduler to run.
func (j *job) Do(fn interface{}, args ...interface{}) {
	fv := reflect.ValueOf(fn)
	if fv.Kind() != reflect.Func {
		panic("Argument fn must be a function type")
	}
	j.f = fv
	ft := fv.Type()
	if ft.NumIn() == 0 {
		j.args = nil
	} else {
		in := make([]reflect.Value, ft.NumIn())
		for i := range in {
			in[i] = reflect.ValueOf(args[i])
			if in[i].Kind() != ft.In(i).Kind() {
				fname := runtime.FuncForPC(fv.Pointer()).Name()
				panic(fmt.Sprintf("cannot use type %T as type %s in argument to %s", args[i], ft.In(i).Name(), fname))
			}
		}
		j.args = in
	}
	j.scheduler.addJob(j)
}
