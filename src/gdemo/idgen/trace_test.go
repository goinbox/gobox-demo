package idgen

import "testing"

func TestTraceIdGenter(t *testing.T) {
	idGenter := NewTraceIdGenter(4)

	for i := 0; i < 100000; i++ {
		id, err := idGenter.GenId("192.168.1.2", "9001")
		t.Log(i, string(id), err)
	}

	for i := 0; i < 100000; i++ {
		id, err := DefaultTraceIdGenter.GenId("192.168.1.2", "9001")
		t.Log(i, string(id), err)
	}
}
