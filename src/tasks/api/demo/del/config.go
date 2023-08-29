package del

import (
	"gdemo/pcontext"
	"github.com/goinbox/taskflow/v2"
)

const (
	stepKeyDeleteFromDB = "deleteFromDB"
)

func (t *Task) FirstStepKey() string {
	return stepKeyDeleteFromDB
}

func (t *Task) StepConfigMap() map[string]*taskflow.StepConfig[*pcontext.Context] {
	return map[string]*taskflow.StepConfig[*pcontext.Context]{
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
