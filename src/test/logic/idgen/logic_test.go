package idgen

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"gdemo/logic/factory"
	"gdemo/test"
)

func init() {
	dir, _ := os.Getwd()
	for i := 0; i < 4; i++ {
		dir = filepath.Dir(dir)
	}

	test.InitTestResource(dir)
}

func TestGenerateID(t *testing.T) {
	ctx := test.Context()
	m := make(map[int64]bool)
	logic := factory.DefaultLogicFactory.IDGenLogic()
	for i := 0; i < 1000000; i++ {
		id := logic.GenerateID(ctx)
		t.Log(id, fmt.Sprintf("%064b", id))

		_, ok := m[id]
		if ok {
			t.Fatalf("duplicated id %d", id)
		}
		m[id] = true
	}
}
