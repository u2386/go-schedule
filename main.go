package main

import (
	"fmt"

	"github.com/u2386/go-schedule/pkg/schedule"
)

var (
	out      = make(chan interface{}, 1)
	scheduler schedule.Scheduler
)

func runWithReturn(s string, o chan<- interface{}) {
	o <- fmt.Sprint(s)
}

func main() {
	var args = []interface{}{"Test", out}
	scheduler.Do(runWithReturn, args...)
	scheduler.RunAll()
}
