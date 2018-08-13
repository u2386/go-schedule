package main

import (
	"fmt"

	"github.com/u2386/go-schedule/pkg/schedule"
)

var (
	out       = make(chan interface{}, 1)
	scheduler schedule.Scheduler
)

func runWithReturn(s string) {
	out <- fmt.Sprint(s)
}

func main() {
	var args = []interface{}{"Test"}
	scheduler.Do(runWithReturn, args...).Every(1).Seconds()
	scheduler.RunOnce()
}
