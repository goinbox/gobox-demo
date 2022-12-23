package list

import (
	"github.com/goinbox/taskflow"
)

const (
	stepKeyMakeSqlQueryParams = "makeSqlQueryParams"
	stepKeyQueryFromDB        = "queryFromDB"
)

func (t *Task) FirstStepKey() string {
	return stepKeyMakeSqlQueryParams
}

func (t *Task) StepConfigMap() map[string]*taskflow.StepConfig {
	return map[string]*taskflow.StepConfig{
		stepKeyMakeSqlQueryParams: {
			RetryCnt:       0,
			RetryDelay:     0,
			StepFunc:       t.makeSqlQueryParams,
			StepFailedFunc: t.StepFailed,
			RouteMap: map[string]string{
				taskflow.StepCodeSuccess: stepKeyQueryFromDB,
				taskflow.StepCodeFailure: "",
			},
		},
		stepKeyQueryFromDB: {
			RetryCnt:       0,
			RetryDelay:     0,
			StepFunc:       t.queryFromDB,
			StepFailedFunc: t.StepFailed,
			RouteMap: map[string]string{
				taskflow.StepCodeSuccess: "",
				taskflow.StepCodeFailure: "",
			},
		},
	}
}
