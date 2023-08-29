package edit

import (
	"testing"
)

func TestMakeUpdateColumns(t *testing.T) {
	code, err := flowTask.makeUpdateColumns(ctx)
	t.Log(code, err, flowTask.data.updateColumns)
}

func TestUpdateDB(t *testing.T) {
	_, _ = flowTask.makeUpdateColumns(ctx)
	code, err := flowTask.updateDB(ctx)
	t.Log(code, err, flowTask.out.RowsAffected)
}
