package edit

import (
	"testing"
)

func TestMakeUpdateColumns(t *testing.T) {
	code, err := flowTask.makeUpdateColumns()
	t.Log(code, err, flowTask.data.updateColumns)
}

func TestUpdateDB(t *testing.T) {
	_, _ = flowTask.makeUpdateColumns()
	code, err := flowTask.updateDB()
	t.Log(code, err, flowTask.out.RowsAffected)
}
