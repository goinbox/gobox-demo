package edit

import (
	"gdemo/pcontext"
	"github.com/goinbox/taskflow/v2"
)

const (
	stepKeyMakeUpdateColumns = "makeUpdateColumns"
	stepKeyUpdateDB          = "updateDB"
)

func (t *Task) FirstStepKey() string {
	return stepKeyMakeUpdateColumns
}

func (t *Task) StepConfigMap() map[string]*taskflow.StepConfig[*pcontext.Context] {
	return map[string]*taskflow.StepConfig[*pcontext.Context]{
		stepKeyMakeUpdateColumns: {
			RetryCnt:       0,
			RetryDelay:     0,
			StepFunc:       t.makeUpdateColumns,
			StepFailedFunc: t.StepFailed,
			RouteMap: map[string]string{
				taskflow.StepCodeSuccess: stepKeyUpdateDB,
				taskflow.StepCodeJump1:   "",
				taskflow.StepCodeFailure: "",
			},
		},
		stepKeyUpdateDB: {
			RetryCnt:       0,
			RetryDelay:     0,
			StepFunc:       t.updateDB,
			StepFailedFunc: t.StepFailed,
			RouteMap: map[string]string{
				taskflow.StepCodeSuccess: "",
				taskflow.StepCodeFailure: "",
			},
		},
	}
}
