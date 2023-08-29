package add

import (
	"gdemo/pcontext"
	"github.com/goinbox/taskflow/v2"
)

const (
	stepKeyGenEntity  = "genEntity"
	stepKeySaveEntity = "saveEntity"
)

func (t *Task) FirstStepKey() string {
	return stepKeyGenEntity
}

func (t *Task) StepConfigMap() map[string]*taskflow.StepConfig[*pcontext.Context] {
	return map[string]*taskflow.StepConfig[*pcontext.Context]{
		stepKeyGenEntity: {
			RetryCnt:       0,
			RetryDelay:     0,
			StepFunc:       t.genEntity,
			StepFailedFunc: t.StepFailed,
			RouteMap: map[string]string{
				taskflow.StepCodeSuccess: stepKeySaveEntity,
				taskflow.StepCodeFailure: "",
			},
		},
		stepKeySaveEntity: {
			RetryCnt:       0,
			RetryDelay:     0,
			StepFunc:       t.saveEntity,
			StepFailedFunc: t.StepFailed,
			RouteMap: map[string]string{
				taskflow.StepCodeSuccess: "",
				taskflow.StepCodeFailure: "",
			},
		},
	}
}
