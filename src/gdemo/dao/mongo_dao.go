package dao

import (
	"github.com/goinbox/mongo"
	"gopkg.in/mgo.v2/bson"
)

type MongoExecResult struct {
	Err error
}

type MongoDao struct {
	*mongo.Client
}

const (
	MONGO_COND_EQUAL         = "$eq"
	MONGO_COND_NOT_EQUAL     = "$ne"
	MONGO_COND_LESS          = "$lt"
	MONGO_COND_LESS_EQUAL    = "$lte"
	MONGO_COND_GREATER       = "$gt"
	MONGO_COND_GREATER_EQUAL = "$gte"
	MONGO_COND_IN            = "$in"
	MONGO_COND_NOT_IN        = "$nin"
)

func NewMongoDao(client *mongo.Client) *MongoDao {
	return &MongoDao{
		client,
	}
}

func (m *MongoDao) InsertRow(tableName string, colNames []string, colValues []interface{}) error {
	doc := make(map[string]interface{})
	for i, entity := range colNames {
		doc[entity] = colValues[i]
	}
	return m.Insert(tableName, doc)
}

func (m *MongoDao) InsertRows(tableName string, colNames []string, colsValues ...[]interface{}) error {
	docs := make([]interface{}, len(colsValues))
	for i, colsValue := range colsValues {
		doc := make(map[string]interface{})
		for j, entity := range colNames {
			doc[entity] = colsValue[j]
		}
		docs[i] = doc
	}
	return m.Insert(tableName, docs...)
}

func (m *MongoDao) DeleteById(tableName string, id interface{}) error {
	return m.RemoveId(tableName, id)
}

func (m *MongoDao) UpdateById(tableName string, id interface{}, updater map[string]interface{}) error {
	return m.UpdateId(tableName, id, map[string]interface{}{"$set": updater})
}

func (m *MongoDao) SelectById(tableName string, id interface{}) (bson.M, error) {
	result := bson.M{}
	err := m.FindId(tableName, id).One(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *MongoDao) SelectByIds(tableName string, ids ...interface{}) ([]bson.M, error) {
	result := []bson.M{}
	err := m.Find(tableName, bson.M{"_id": bson.M{"$in": ids}}).All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *MongoDao) SelectCount(tableName string, selector interface{}) (int, error) {
	n, err := m.Find(tableName, selector).Count()
	if err != nil {
		return 0, err
	}
	return n, err
}

func (m *MongoDao) SelectAll(tableName string, query *mongo.Query) ([]bson.M, error) {
	result := []bson.M{}
	err := m.Query(tableName, query).All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *MongoDao) SelectOne(tableName string, query *mongo.Query) (bson.M, error) {
	result := bson.M{}
	err := m.Query(tableName, query).One(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
