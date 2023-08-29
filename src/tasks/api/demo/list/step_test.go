package list

import (
	"testing"
)

func TestMakeSqlQueryParams(t *testing.T) {
	code, err := flowTask.makeSqlQueryParams(ctx)
	t.Log(code, err, flowTask.data.queryParams)
	for _, item := range flowTask.data.queryParams.CondItems {
		t.Log(item)
	}
}

func TestQueryFromDB(t *testing.T) {
	_, _ = flowTask.makeSqlQueryParams(ctx)
	code, err := flowTask.queryFromDB(ctx)
	t.Log(code, err, flowTask.out)
	for _, entity := range flowTask.out.DemoList {
		t.Log(entity, *entity.ID, *entity.AddTime, *entity.EditTime)
	}
}
