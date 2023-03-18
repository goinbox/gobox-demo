package add

import (
	"testing"
)

func TestGenEntity(t *testing.T) {
	code, err := flowTask.genEntity()
	t.Log(code, err, flowTask.data.demoEntity)
}

func TestSaveEntity(t *testing.T) {
	_, _ = flowTask.genEntity()
	code, err := flowTask.saveEntity()
	t.Log(code, err, flowTask.out)
}
