package misc

import "testing"

func TestTraceIDGenter(t *testing.T) {
	idGenter := NewTraceIDGenter(4)

	for i := 0; i < 100000; i++ {
		id, err := idGenter.GenID("192.168.1.2", 9001)
		t.Log(i, string(id), err)
	}

	for i := 0; i < 100000; i++ {
		id, err := DefaultTraceIDGenter.GenID("192.168.1.2", 9001)
		t.Log(i, string(id), err)
	}
}
