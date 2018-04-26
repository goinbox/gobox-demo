package mongodemo

import (
	"gdemo/dao"
	"gdemo/misc"
	"gdemo/svc"
	"time"

	"testing"

	"github.com/goinbox/gomisc"
)

type testQueryParamsStruct struct {
	Name   string `bson:"name" json:"name" redis:"name"`
	Status int    `bson:"status" json:"status" redis:"status"`
}

func TestDemoSvc(t *testing.T) {
	demoSvc := NewMongoDemoSvc(nil, misc.MongoTestClient())

	ids, err := demoSvc.Insert(
		&MongoDemoEntity{Name: "a1", Status: 0},
		&MongoDemoEntity{Name: "a2", Status: 1},
	)
	t.Log(ids, err)

	for _, id := range ids {
		entity, err := demoSvc.GetById(id)
		t.Log(entity, err)

		deleted, err := demoSvc.DeleteById(id)
		t.Log(deleted, err)
	}

	baseEntity := svc.MongoBaseEntity{AddTime: time.Now().Format(gomisc.TimeGeneralLayout())}
	demoSvc.UpdateById(11, &MongoDemoEntity{MongoBaseEntity: baseEntity, Name: "ccc", Status: 44}, map[string]bool{"name": true, "status": true})

	entity, err := demoSvc.GetById(11)
	t.Log(entity, err)

	mqp := &svc.MongoQueryParams{
		ParamsStructPtr: &testQueryParamsStruct{
			Name:   "ccc",
			Status: 1,
		},
		Exists: map[string]bool{"name": true, "status": true},
		Conditions: map[string]string{
			"name":   dao.MONGO_COND_EQUAL,
			"status": dao.MONGO_COND_GREATER_EQUAL,
		},
		OrderBy: "_id", Offset: 0, Cnt: 10,
	}
	entities, err := demoSvc.SelectAll(mqp)
	t.Log(entities, err)

}
