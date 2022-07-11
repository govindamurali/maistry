package maistry

import (
	"errors"
)

type Job struct {
	do     func()
	logger ILogger
}

type jobChannel chan Job

func NewJob(do func(), extLogger ILogger) Job {
	return Job{
		do:     do,
		logger: extLogger,
	}
}

func (j *Job) Do() {
	defer func() {
		if r := recover(); r != nil {
			j.logger.Error("Job", errors.New("panic occured in a Job"), map[string]interface{}{"recover_type": r})
		}
	}()

	j.do()
}
