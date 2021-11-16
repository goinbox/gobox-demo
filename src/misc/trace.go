package misc

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
)

type TraceIDGenter struct {
	lock      sync.Mutex
	increment int64

	incrementLen    int
	maxIncrement    int64
	incrementFormat string
}

func NewTraceIDGenter(incrementLen int) *TraceIDGenter {
	return new(TraceIDGenter).SetIncrementLen(incrementLen)
}

var DefaultTraceIDGenter = NewTraceIDGenter(4)

func (t *TraceIDGenter) SetIncrementLen(incrementLen int) *TraceIDGenter {
	t.incrementLen = incrementLen
	t.maxIncrement = int64(math.Pow10(incrementLen))
	t.incrementFormat = "%0" + strconv.Itoa(incrementLen) + "d"

	return t
}

func (t *TraceIDGenter) GenID(ip string, port int) (string, error) {
	var id string

	for _, item := range strings.Split(ip, ".") {
		v, err := strconv.Atoi(item)
		if err != nil {
			return "", fmt.Errorf("strconv.Atoi error: %w", err)
		}
		id += fmt.Sprintf("%02x", v)
	}

	id += fmt.Sprintf("%05d", port)
	id += strconv.FormatInt(time.Now().UnixNano()/1000000, 10)

	t.lock.Lock()
	increment := t.increment
	t.increment = (t.increment + 1) % t.maxIncrement
	t.lock.Unlock()

	id += fmt.Sprintf(t.incrementFormat, increment)

	return id, nil
}
