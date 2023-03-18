package del

import (
	"testing"
)

func TestTaskGraph(t *testing.T) {
	graph := runner.TaskGraph(flowTask)

	t.Log(graph)
}
