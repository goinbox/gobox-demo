package del

import (
	"github.com/goinbox/taskflow"
)

const (
	stepKeyDeleteFromDB = "deleteFromDB"
)

func (t *Task) FirstStepKey() string {
	return stepKeyDeleteFromDB
}

func (t *Task) StepConfigMap() map[string]*taskflow.StepConfig {
	return map[string]*taskflow.StepConfig{
		stepKeyDeleteFromDB: {
			RetryCnt:       0,
			RetryDelay:     0,
			StepFunc:       t.deleteFromDB,
			StepFailedFunc: t.StepFailed,
			RouteMap: map[string]string{
				taskflow.StepCodeSuccess: "",
				taskflow.StepCodeFailure: "",
			},
		},
	}
}
