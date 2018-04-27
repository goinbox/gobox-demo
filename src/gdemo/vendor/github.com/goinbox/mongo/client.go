package mongo

import (
	"encoding/json"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/goinbox/golog"

	"fmt"
)

type CmdLogFmtFunc func(cmd string, args ...interface{}) []byte

type Client struct {
	config *Config
	logger golog.ILogger
	clff   CmdLogFmtFunc

	conn      *mgo.Session
	db        *mgo.Database
	coll      *mgo.Collection
	connected bool

	pipeCnt int
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
		logger: logger,
	}
	c.clff = c.cmdLogFmt
	c.SetLogger(logger)

	return c
}

func (c *Client) SetLogger(logger golog.ILogger) *Client {
	c.logger = logger
	mgo.SetLogger(NewMongoLogger(logger))

	return c
}

func (c *Client) SetDebug(debug bool) {
	mgo.SetDebug(debug)
}

func (c *Client) SetCmdLogFmtFunc(clff CmdLogFmtFunc) *Client {
	c.clff = clff

	return c
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
	url := "mongodb://" + c.config.User + ":" + c.config.Pass + "@" + c.config.Host + ":" + c.config.Port

	session, err := mgo.Dial(url)
	if err != nil {
		return err
	}

	//session.SetMode(mgo.Monotonic, true)
	session.SetMode(mgo.Eventual, true)

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
	q := c.Collection(coll).Find(query.finder)
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
	return c.BuildQuery(coll, query)
}

func (c *Client) Find(coll string, finder interface{}) *mgo.Query {
	return c.Collection(coll).Find(finder)
}

func (c *Client) FindId(coll string, id interface{}) *mgo.Query {
	return c.Collection(coll).FindId(id)
}

func (c *Client) FindAndModify(coll string, finder interface{}, updater interface{}) (result bson.M, err error) {
	change := mgo.Change{
		Update:    updater,
		Upsert:    true,
		ReturnNew: true,
	}

	result = bson.M{}
	_, err = c.Collection(coll).Find(finder).Apply(change, result)
	if err != nil {
		c.log("FindAndModify Fail, Error:", err)
	}
	return result, err
}

func (c *Client) Indexes(coll string) (indexes []mgo.Index, err error) {
	indexes, err = c.Collection(coll).Indexes()
	if err != nil {
		c.log("Indexes Fail, Indexes:", indexes,
			", Error:", err)
	}
	return indexes, err
}

func (c *Client) Insert(coll string, docs ...interface{}) error {
	err := c.Collection(coll).Insert(docs...)
	if err != nil {
		c.log("Insert Fail, docs:", docs,
			", Error:", err)
	}
	return err
}

func (c *Client) Update(coll string, selector, updater interface{}) error {
	err := c.Collection(coll).Update(selector, updater)
	if err != nil {
		c.log("Update Fail, Selector:", selector,
			", Updater:", updater,
			", Error:", err)
	}
	return err
}

func (c *Client) UpdateAll(coll string, selector, updater interface{}) error {
	_, err := c.Collection(coll).UpdateAll(selector, updater)
	if err != nil {
		c.log("UpdateAll Fail, Selector:", selector,
			", Updater:", updater,
			", Error:", err)
	}
	return err
}

func (c *Client) UpdateId(coll string, id interface{}, updater interface{}) error {
	err := c.Collection(coll).UpdateId(id, updater)
	if err != nil {
		c.log("UpdateId Fail, Id:", id,
			", Updater:", updater,
			", Error:", err)
	}
	return err
}

func (c *Client) Upsert(coll string, selector, updater interface{}) error {
	_, err := c.Collection(coll).Upsert(selector, updater)
	if err != nil {
		c.log("Upsert Fail, Selector:", selector,
			", Updater:", updater,
			", Error:", err)
	}
	return err
}

func (c *Client) Remove(coll string, selector interface{}) error {
	err := c.Collection(coll).Remove(selector)
	if err != nil {
		c.log("Remove Fail, Selector:", selector,
			", Error:", err)
	}
	return err
}

func (c *Client) RemoveAll(coll string, selector interface{}) error {
	_, err := c.Collection(coll).RemoveAll(selector)
	if err != nil {
		c.log("RemoveAll Fail, Selector:", selector,
			", Error:", err)
	}
	return err
}

func (c *Client) RemoveId(coll string, id interface{}) error {
	err := c.Collection(coll).RemoveId(id)
	if err != nil {
		c.log("RemoveId Fail, Id:", id,
			", Error:", err)
	}
	return err
}

func (c *Client) log(cmd string, args ...interface{}) {
	if len(cmd) == 0 {
		return
	}

	msg := c.clff(cmd, args...)
	if msg != nil {
		c.logger.Log(c.config.LogLevel, msg)
	}
}

func (c *Client) cmdLogFmt(cmd string, args ...interface{}) []byte {
	for _, arg := range args {
		cmd += " " + fmt.Sprint(arg)
	}

	return []byte(cmd)
}

func (c *Client) ConvertBsonToStruct(doc interface{}, entity interface{}) error {
	bsonBytes, _ := bson.Marshal(doc)
	err := bson.Unmarshal(bsonBytes, entity)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) ConvertJsonToStruct(doc interface{}, entity interface{}) error {
	jsonBytes, _ := json.Marshal(doc)
	err := json.Unmarshal(jsonBytes, entity)
	if err != nil {
		return err
	}
	return nil
}
