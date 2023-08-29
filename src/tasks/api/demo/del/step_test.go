package del

import (
	"testing"
)

func TestDeleteFromDB(t *testing.T) {
	code, err := flowTask.deleteFromDB(ctx)
	t.Log(code, err, flowTask.out.RowsAffected)
}
