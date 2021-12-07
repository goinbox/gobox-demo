package misc

import "testing"

func TestTraceIDGenter(t *testing.T) {
	idGenter := NewTraceIDGenter(4)

	for i := 0; i < 100000; i++ {
		id, err := idGenter.GenID("1.2.3.4", 9001)
		t.Log(i, id, err)
	}

	for i := 0; i < 100000; i++ {
		id, err := DefaultTraceIDGenter.GenID("1.2.3.4", 9001)
		t.Log(i, id, err)
	}
}
