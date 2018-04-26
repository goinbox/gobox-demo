package idgen

import (
	"github.com/goinbox/mongo"
	"gopkg.in/mgo.v2/bson"
)

type MongoIdGenter struct {
	client *mongo.Client
}

type IdGenter struct {
	Id    bson.ObjectId `bson:"_id,omitempty"`
	MaxId int64         `bson:"max_id"`
}

func NewMongoIdGenter(client *mongo.Client) *MongoIdGenter {
	return &MongoIdGenter{
		client: client,
	}
}

func (m *MongoIdGenter) GenId(name string) (int64, error) {
	finder := bson.M{"_id": name}
	updater := bson.M{"$inc": bson.M{"max_id": 1}}

	result, err := m.client.FindAndModify("id_genter", finder, updater)
	if err != nil {
		return 0, err
	}

	entity := new(IdGenter)
	err = m.client.ConvertBsonToStruct(result, entity)
	if err != nil {
		return 0, err
	}

	return entity.MaxId, nil
}
