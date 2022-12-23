package list

import (
	"testing"

	"github.com/goinbox/taskflow"
)

func TestTaskGraph(t *testing.T) {
	graph := taskflow.NewRunner(nil).TaskGraph(NewTask(nil))

	t.Log(graph)
}
