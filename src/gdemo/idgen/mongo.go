package idgen

import (
	"github.com/globalsign/mgo/bson"
	"github.com/goinbox/mongo"
)

type MongoIdGenter struct {
	client *mongo.Client
}

type IdGenterItem struct {
	Id    bson.ObjectId `bson:"_id,omitempty"`
	MaxId int64         `bson:"max_id"`
}

func NewMongoIdGenter(client *mongo.Client) *MongoIdGenter {
	return &MongoIdGenter{
		client: client,
	}
}

func (m *MongoIdGenter) SetClient(client *mongo.Client) *MongoIdGenter {
	m.client = client

	return m
}

func (m *MongoIdGenter) GenId(tableName string) (int64, error) {
	finder := bson.M{"_id": tableName}
	updater := bson.M{"$inc": bson.M{"max_id": 1}}

	result, err := m.client.FindAndModify("id_genter", finder, updater)
	if err != nil {
		return 0, err
	}

	entity := new(IdGenterItem)
	err = m.client.ConvertBsonToStruct(result, entity)
	if err != nil {
		return 0, err
	}

	return entity.MaxId, nil
}
