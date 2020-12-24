package mongodemo

import (
	"os"
	"testing"
	"time"

	"github.com/goinbox/mongo"

	"gdemo/conf"
	"gdemo/define"
	"gdemo/define/entity"
	"gdemo/resource"
)

type testQueryParamsStruct struct {
	Name   string `bson:"name" json:"name"`
	Status int    `bson:"status" json:"status"`
}

func init() {
	_ = conf.Init(os.Getenv("GOPATH"))

	_ = resource.InitLog("test")
	resource.InitMongo()
}

func TestDemoSvc(t *testing.T) {
	demoSvc := NewMongoDemoSvc([]byte("traceMongoDemoSvc"))

	ids, err := demoSvc.Insert(
		&entity.MongoDemoEntity{Name: "a1", Status: 0},
		&entity.MongoDemoEntity{Name: "a2", Status: 1},
	)
	t.Log(ids, err)

	for _, id := range ids {
		entity, err := demoSvc.GetById(id)
		t.Log(entity, err)

		deleted, err := demoSvc.DeleteById(id)
		t.Log(deleted, err)
	}

	baseEntity := entity.MongoBaseEntity{AddTime: time.Now()}
	demoSvc.UpdateById(1, &entity.MongoDemoEntity{MongoBaseEntity: baseEntity, Name: "ccc", Status: 1}, map[string]bool{"name": true, "status": true, "add_time": true})

	entity, err := demoSvc.GetById(1)
	t.Log(entity, err)

	mqp := &define.MongoQueryParams{
		ParamsStructPtr: &testQueryParamsStruct{
			Name:   "ccc",
			Status: 1,
		},
		Exists: map[string]bool{"name": true, "status": true},
		Conditions: map[string]string{
			"name":   mongo.MONGO_COND_EQUAL,
			"status": mongo.MONGO_COND_LESS_EQUAL,
		},
		OrderBy: []string{"name", "-_id"}, Offset: 0, Cnt: 10,
	}
	entities, err := demoSvc.SelectAll(mqp)
	t.Log(entities, err)

	mqp = &define.MongoQueryParams{
		ParamsStructPtr: &testQueryParamsStruct{
			Name:   "CC.+",
			Status: 1,
		},
		Exists: map[string]bool{"name": true, "status": true},
		Conditions: map[string]string{
			"name":   mongo.MONGO_COND_REGEX,
			"status": mongo.MONGO_COND_LESS_EQUAL,
		},
		OrderBy: []string{"name", "-_id"}, Offset: 0, Cnt: 10,
	}
	entities, err = demoSvc.SelectRegex(mqp)
	t.Log(entities, err)
}
