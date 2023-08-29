package tasks

import (
	"fmt"
)

type BaseTask struct {
	err error
}

func NewBaseTask() *BaseTask {
	return &BaseTask{}
}

func (t *BaseTask) BeforeStep(stepKey string) {
}

func (t *BaseTask) AfterStep(stepKey string) {
}

func (t *BaseTask) StepFailed(stepKey string, err error) {
	t.err = fmt.Errorf("step %s run error: %w", stepKey, err)
}

func (t *BaseTask) Error() error {
	return t.err
}

func (t *BaseTask) SetError(err error) {
	t.err = err
}
