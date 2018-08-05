package schedule

import "reflect"

type job struct {
	f        reflect.Value
	args     []reflect.Value
	interval int
	unit     int
}

func (j *job) Run() {
	j.f.Call(j.args)
}
