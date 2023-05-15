package utils

type JobChannel[F func()] struct {
	callBack F
}

func (j *JobChannel[F]) AddJob(callBack F) {
	current := j.callBack
	if current == nil {
		j.callBack = callBack
	} else {
		j.callBack = func() {
			current()
			callBack()
		}
	}
}

func (j *JobChannel[F]) Run() {
	current := j.callBack
	if current != nil {
		current()
	}
}
