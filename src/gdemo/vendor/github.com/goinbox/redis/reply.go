package redis

import "github.com/garyburd/redigo/redis"

type Reply struct {
	reply interface{}

	Err error
}

func NewReply(reply interface{}, err error) *Reply {
	return &Reply{
		reply: reply,

		Err: err,
	}
}

func (r *Reply) SimpleReplyIsNil() bool {
	if r.reply == nil {
		return true
	}

	return false
}

func (r *Reply) ArrReplyIsNil() bool {
	v, _ := redis.Values(r.reply, nil)
	if len(v) == 0 {
		return true
	}

	return false
}

func (r *Reply) Bool() (bool, error) {
	return redis.Bool(r.reply, nil)
}

func (r *Reply) ByteSlices() ([][]byte, error) {
	return redis.ByteSlices(r.reply, nil)
}

func (r *Reply) Bytes() ([]byte, error) {
	return redis.Bytes(r.reply, nil)
}

func (r *Reply) Float64() (float64, error) {
	return redis.Float64(r.reply, nil)
}

func (r *Reply) Int() (int, error) {
	return redis.Int(r.reply, nil)
}

func (r *Reply) Int64() (int64, error) {
	return redis.Int64(r.reply, nil)
}

func (r *Reply) Int64Map() (map[string]int64, error) {
	return redis.Int64Map(r.reply, nil)
}

func (r *Reply) Ints() ([]int, error) {
	return redis.Ints(r.reply, nil)
}

func (r *Reply) Int64s() ([]int64, error) {
	return redis.Int64s(r.reply, nil)
}

func (r *Reply) Struct(s interface{}) error {
	return redis.ScanStruct(r.reply.([]interface{}), s)
}

func (r *Reply) String() (string, error) {
	return redis.String(r.reply, nil)
}

func (r *Reply) StringMap() (map[string]string, error) {
	return redis.StringMap(r.reply, nil)
}

func (r *Reply) Strings() ([]string, error) {
	return redis.Strings(r.reply, nil)
}

func (r *Reply) Uint64() (uint64, error) {
	return redis.Uint64(r.reply, nil)
}
