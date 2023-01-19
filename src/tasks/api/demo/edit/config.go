package edit

import (
	"github.com/goinbox/taskflow"
)

const (
	stepKeyMakeUpdateColumns = "makeUpdateColumns"
	stepKeyUpdateDB          = "updateDB"
)

func (t *Task) FirstStepKey() string {
	return stepKeyMakeUpdateColumns
}

func (t *Task) StepConfigMap() map[string]*taskflow.StepConfig {
	return map[string]*taskflow.StepConfig{
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
