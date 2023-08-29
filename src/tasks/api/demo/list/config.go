package list

import (
	"gdemo/pcontext"
	"github.com/goinbox/taskflow/v2"
)

const (
	stepKeyMakeSqlQueryParams = "makeSqlQueryParams"
	stepKeyQueryFromDB        = "queryFromDB"
)

func (t *Task) FirstStepKey() string {
	return stepKeyMakeSqlQueryParams
}

func (t *Task) StepConfigMap() map[string]*taskflow.StepConfig[*pcontext.Context] {
	return map[string]*taskflow.StepConfig[*pcontext.Context]{
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
