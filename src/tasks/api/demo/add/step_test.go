package add

import (
	"testing"
)

func TestGenEntity(t *testing.T) {
	code, err := flowTask.genEntity(ctx)
	t.Log(code, err, flowTask.data.demoEntity)
}

func TestSaveEntity(t *testing.T) {
	_, _ = flowTask.genEntity(ctx)
	code, err := flowTask.saveEntity(ctx)
	t.Log(code, err, flowTask.out)
}
