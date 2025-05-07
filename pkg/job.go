package pkg

import (
	"fmt"
	"time"
)

const (
	IDLE = iota
	RUNNING
)

type Job struct {
	prefix  string
	timer   time.Duration
	handler func() error
	status  int
}

func NewJob(prefix string, timer time.Duration, handler func() error) *Job {
	return &Job{prefix: prefix, timer: timer, handler: handler}
}

func (j *Job) Run() error {
	timer := time.NewTimer(j.timer)

	for {
		select {
		case <-timer.C:
			fmt.Println(fmt.Sprintf("%v is processing", j.prefix))
			j.execute()
			timer.Reset(j.timer)
		}
	}
}

func (j *Job) execute() {
	defer func() {
		j.status = IDLE
		if r := recover(); r != nil {
			fmt.Println(fmt.Sprintf("%v job error: %v", j.prefix, r))
		}
	}()

	if j.status == RUNNING {
		return
	}

	j.status = RUNNING
	if err := j.handler(); err != nil {
		fmt.Println(fmt.Sprintf("%v is processed with %v", j.prefix, err.Error()))
	}
}
