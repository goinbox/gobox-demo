package store

import (
	"reflect"
	"testing"

	"github.com/goinbox/mongo"

	"gdemo/define"
	"gdemo/define/entity"
	"gdemo/resource"
)

var MongoDemoEntityType = reflect.TypeOf(MongoDemoEntity{})
var MongoDemoColNames = entity.ReflectMongoColNames(MongoDemoEntityType)

type MongoDemoEntity struct {
	entity.MongoBaseEntity

	Name   string `bson:"name" json:"name"`
	Status int    `bson:"status" json:"status"`
}

func TestMongoSvcInsertGetListUpdateDelete(t *testing.T) {
	resource.InitMongoTest()
	ms := NewMongoStore([]byte("TestMongoSvcInsertGetListUpdateDelete"), resource.MongoClientPool, true)

	resource.TestLogger.Notice([]byte("test Insert"))

	item := &MongoDemoEntity{
		Name:   "ccc",
		Status: 1,
	}

	tableName := "demo"
	ids, err := ms.Insert(tableName, MongoDemoColNames, item)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ids)

	resource.TestLogger.Notice([]byte("test List"))

	mqp := &define.MongoQueryParams{
		ParamsStructPtr: &MongoDemoEntity{
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
	entities := new([]MongoDemoEntity)
	err = ms.SelectAll(entities, tableName, mqp, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(entities)

	resource.TestLogger.Notice([]byte("test Update"))

	newDemo := &MongoDemoEntity{
		Name: "new-demo",
	}
	updateFields := map[string]bool{"name": true}
	setItems, err := ms.UpdateById(tableName, ids[0], newDemo, updateFields)
	if err != nil {
		t.Fatal(err)
	}
	for i, item := range setItems {
		t.Log(i, item)
	}

	resource.TestLogger.Notice([]byte("test Get"))

	item = &MongoDemoEntity{}
	find, err := ms.GetById(item, tableName, ids[0])
	if !find {
		t.Fatal("not found")
	}
	if err != nil {
		t.Fatal(err)
	}
	t.Log(item)

	resource.TestLogger.Notice([]byte("test Delete"))

	result := ms.Dao().DeleteById(tableName, ids[0])

	t.Log(result)
}
