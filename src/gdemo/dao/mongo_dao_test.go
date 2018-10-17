package dao

//import (
//	"gdemo/misc"
//
//	"github.com/goinbox/gomisc"
//	"github.com/goinbox/mongo"
//	"gopkg.in/mgo.v2/bson"
//
//	"testing"
//	"time"
//)
//
//const (
//	MONGO_TEST_TABLE_NAME = "testcoll"
//)
//
//type mongoTestEntity struct {
//	Id       int64  `bson:"_id" json:"_id"`
//	Name     string `bson:"name" json:"name"`
//	Status   int    `bson:"status" json:"status"`
//	AddTime  string `bson:"add_time" json:"add_time"`
//	EditTime string `bson:"edit_time" json:"edit_time"`
//}
//
//func TestMongoDaoRead(t *testing.T) {
//	dao := &MongoDao{misc.MongoTestClient()}
//
//	result, err := dao.SelectById(MONGO_TEST_TABLE_NAME, 11)
//	t.Log(result, err)
//	entity := new(mongoTestEntity)
//	err = dao.ConvertBsonToStruct(result, entity)
//	t.Log(entity, err)
//
//	results, err := dao.SelectByIds(MONGO_TEST_TABLE_NAME, 100, 101)
//	t.Log(results, err)
//	entities := new([]mongoTestEntity)
//	err = dao.ConvertJsonToStruct(results, entities)
//	t.Log(entities, err)
//
//	count, err := dao.SelectCount(MONGO_TEST_TABLE_NAME, bson.M{"_id": bson.M{"$gt": 11}})
//	t.Log(count, err)
//
//	queryAll := mongo.NewQuery().Find(bson.M{"_id": bson.M{"$gte": 100}}).Sort("-_id").Select(bson.M{"edit_time": 0}).Skip(0).SetMaxTime(1 * time.Second)
//	rows, err := dao.SelectAll(MONGO_TEST_TABLE_NAME, queryAll)
//	t.Log(rows, err)
//
//	queryOne := mongo.NewQuery().Find(bson.M{"_id": 100}).Sort("_id").Select(bson.M{"edit_time": 0}).Skip(0).SetMaxTime(1 * time.Second)
//	row, err := dao.SelectOne(MONGO_TEST_TABLE_NAME, queryOne)
//	t.Log(row, err)
//}
//
//func TestMongoDaoWrite(t *testing.T) {
//	dao := &MongoDao{misc.MongoTestClient()}
//
//	var colNames = []string{"_id", "add_time", "edit_time", "name", "status"}
//	var colValues []interface{}
//
//	ts := time.Now().Format(gomisc.TimeGeneralLayout())
//	colValues = []interface{}{
//		11,
//		ts,
//		ts,
//		"a",
//		time.Now().Unix() % 10,
//	}
//	result := dao.InsertRow(MONGO_TEST_TABLE_NAME, colNames, colValues)
//	t.Log(result)
//
//	var colsValues [][]interface{}
//	for i, name := range []string{"a", "b", "c"} {
//		colValues := []interface{}{
//			int64(i + 100),
//			ts,
//			ts,
//			name,
//			i % 10,
//		}
//		colsValues = append(colsValues, colValues)
//	}
//	result = dao.InsertRows(MONGO_TEST_TABLE_NAME, colNames, colsValues...)
//	t.Log(result)
//
//	result = dao.UpdateById(MONGO_TEST_TABLE_NAME, 11, map[string]interface{}{"add_time": "123", "name": "dd"})
//	t.Log(result)
//}
