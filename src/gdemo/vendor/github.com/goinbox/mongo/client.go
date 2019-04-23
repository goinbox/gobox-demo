package mongo

import (
	"encoding/json"

	"github.com/goinbox/golog"
	"github.com/goinbox/gomisc"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"

	"fmt"
	"strings"
)

type Client struct {
	config *Config

	conn      *mgo.Session
	db        *mgo.Database
	coll      *mgo.Collection
	connected bool

	pipeCnt int

	logger    golog.ILogger
	traceId   []byte
	logPrefix []byte
}

func NewClient(config *Config, logger golog.ILogger) *Client {
	if config.LogLevel == 0 {
		config.LogLevel = golog.LEVEL_INFO
	}

	if logger == nil {
		logger = new(golog.NoopLogger)
	}

	c := &Client{
		config: config,

		logger:  logger,
		traceId: []byte("-"),
		logPrefix: []byte("[MongoClient " +
			strings.Join(config.Hosts, ",") +
			"]\t"),
	}

	return c
}

func (c *Client) SetLogger(logger golog.ILogger) *Client {
	if logger == nil {
		logger = new(golog.NoopLogger)
	}
	c.logger = logger

	return c
}

func (c *Client) SetTraceId(traceId []byte) *Client {
	c.traceId = traceId

	return c
}

func (c *Client) SetDebug(debug bool) {
	mgo.SetDebug(debug)
}

func (c *Client) Connected() bool {
	return c.connected
}

func (c *Client) Free() {
	if c.conn != nil {
		c.conn.Close()
	}

	c.connected = false
}

func (c *Client) Connect() error {
	url := "mongodb://"
	if c.config.User == "" && c.config.Pass == "" {
		url += strings.Join(c.config.Hosts, ",")
	} else {
		url += c.config.User + ":" + c.config.Pass + "@" + strings.Join(c.config.Hosts, ",")
	}

	session, err := mgo.Dial(url)
	if err != nil {
		return err
	}

	session.SetMode(c.config.Mode, true)
	session.SetSocketTimeout(c.config.SocketTimeout)
	session.SetSyncTimeout(c.config.SyncTimeout)

	c.conn = session
	c.db = session.DB(c.config.DBName)
	c.connected = true

	return nil
}

func (c *Client) connectCheck() {
	if !c.Connected() {
		if err := c.Connect(); err != nil {
			panic(err)
		}
	}
}

func (c *Client) DB(name string) *mgo.Database {
	database := c.conn.DB(name)
	c.db = database
	return database
}

func (c *Client) Collection(coll string) *mgo.Collection {
	c.connectCheck()
	c.coll = c.db.C(coll)
	return c.coll
}

func (c *Client) Count(coll string) (n int, err error) {
	n, err = c.Collection(coll).Count()
	if err != nil {
		return 0, err
	}
	return n, err
}

func (c *Client) BuildQuery(coll string, query *Query) *mgo.Query {
	q := c.Collection(coll).Find(query.finder).SetMaxTime(c.config.QueryTimeout)
	if query.selector != nil {
		q = q.Select(query.selector)
	}
	if query.sort != nil {
		q = q.Sort(query.sort...)
	}
	if query.limit != 0 {
		q = q.Limit(query.limit)
	}
	if query.skip != 0 {
		q = q.Skip(query.skip)
	}
	if query.setMaxTime != 0 {
		q = q.SetMaxTime(query.setMaxTime)
	}
	return q
}

func (c *Client) Query(coll string, query *Query) *mgo.Query {
	c.log("Query", coll, query)
	return c.BuildQuery(coll, query)
}

func (c *Client) Find(coll string, finder interface{}) *mgo.Query {
	c.log("Find", coll, finder)
	return c.Collection(coll).Find(finder)
}

func (c *Client) FindId(coll string, id interface{}) *mgo.Query {
	c.log("FindId", coll, id)
	return c.Collection(coll).FindId(id)
}

func (c *Client) FindAndModify(coll string, finder interface{}, updater interface{}) (result bson.M, err error) {
	c.log("FindAndModify", coll, finder, updater)
	change := mgo.Change{
		Update:    updater,
		Upsert:    true,
		ReturnNew: true,
	}

	result = bson.M{}
	_, err = c.Collection(coll).Find(finder).Apply(change, result)
	return result, err
}

func (c *Client) Indexes(coll string) (indexes []mgo.Index, err error) {
	c.log("Indexes", coll)
	indexes, err = c.Collection(coll).Indexes()
	return indexes, err
}

func (c *Client) Insert(coll string, docs ...interface{}) error {
	c.log("Insert", coll, docs)
	err := c.Collection(coll).Insert(docs...)
	return err
}

func (c *Client) Update(coll string, selector, updater interface{}) error {
	c.log("Update", coll, selector, updater)
	err := c.Collection(coll).Update(selector, updater)
	return err
}

func (c *Client) UpdateAll(coll string, selector, updater interface{}) error {
	c.log("UpdateAll", coll, selector, updater)
	_, err := c.Collection(coll).UpdateAll(selector, updater)
	return err
}

func (c *Client) UpdateId(coll string, id interface{}, updater interface{}) error {
	c.log("UpdateId", coll, id, updater)
	err := c.Collection(coll).UpdateId(id, updater)
	return err
}

func (c *Client) Upsert(coll string, selector, updater interface{}) error {
	c.log("Upsert", coll, selector, updater)
	_, err := c.Collection(coll).Upsert(selector, updater)
	return err
}

func (c *Client) Remove(coll string, selector interface{}) error {
	c.log("Remove", coll, selector)
	err := c.Collection(coll).Remove(selector)
	return err
}

func (c *Client) RemoveAll(coll string, selector interface{}) error {
	c.log("RemoveAll", coll, selector)
	_, err := c.Collection(coll).RemoveAll(selector)
	return err
}

func (c *Client) RemoveId(coll string, id interface{}) error {
	c.log("RemoveId", coll, id)
	err := c.Collection(coll).RemoveId(id)
	return err
}

func (c *Client) ConvertBsonToStruct(doc interface{}, entity interface{}) error {
	bsonBytes, _ := bson.Marshal(doc)
	err := bson.Unmarshal(bsonBytes, entity)
	return err
}

func (c *Client) ConvertJsonToStruct(doc interface{}, entity interface{}) error {
	jsonBytes, _ := json.Marshal(doc)
	err := json.Unmarshal(jsonBytes, entity)
	return err
}

func (c *Client) log(query string, args ...interface{}) {
	for _, arg := range args {
		query += " " + fmt.Sprint(arg)
	}
	_ = c.logger.Log(c.config.LogLevel, gomisc.AppendBytes(c.traceId, []byte("\t"), c.logPrefix, []byte(query)))
}
