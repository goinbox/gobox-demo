package idgen

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
)

type TraceIdGenter struct {
	lock      sync.Mutex
	increment int64

	incrementLen    int
	maxIncrement    int64
	incrementFormat string
}

func NewTraceIdGenter(incrementLen int) *TraceIdGenter {
	return new(TraceIdGenter).SetIncrementLen(incrementLen)
}

var DefaultTraceIdGenter = NewTraceIdGenter(4)

func (t *TraceIdGenter) SetIncrementLen(incrementLen int) *TraceIdGenter {
	t.incrementLen = incrementLen
	t.maxIncrement = int64(math.Pow10(incrementLen))
	t.incrementFormat = "%0" + strconv.Itoa(incrementLen) + "d"

	return t
}

func (t *TraceIdGenter) GenId(ip, port string) ([]byte, error) {
	var id string

	for _, item := range strings.Split(ip, ".") {
		v, err := strconv.Atoi(item)
		if err != nil {
			return nil, err
		}
		id += fmt.Sprintf("%02x", v)
	}

	id += fmt.Sprintf("%05s", port)
	id += strconv.FormatInt(time.Now().UnixNano()/1000000, 10)

	t.lock.Lock()
	increment := t.increment
	t.increment = (t.increment + 1) % t.maxIncrement
	t.lock.Unlock()

	id += fmt.Sprintf(t.incrementFormat, increment)

	return []byte(id), nil
}
