package workflow

import (
	"testing"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

// UI  http://localhost:8233/

func TestWk(t *testing.T) {
	//  /usr/local/go/bin/go test -timeout 600s -run ^TestWk$ github.com/chernyshev-alex/go-hexagon/internal/adapters/workflow -count=1 -v
	pipeline, order := map[string][]string{
		"fork-1": {"branch-1", "branch-2"},
		"fork-2": {"branch-3", "branch-4", "branch-5"},
		"fork-3": {"branch-6", "branch-7", "branch-8"},
		"fork-4": {"branch-10", "branch-11", "branch-12"},
	},
		[]string{"fork-1", "fork-2", "fork-3", "fork-4"}

	RunWorflow(pipeline, order)
}

func TestWkWorker(t *testing.T) {
	// go test -timeout 600s -run ^TestWkWorker$ github.com/chernyshev-alex/go-hexagon/internal/adapters/workflow -count=1 -v
	c, _ := client.Dial(client.Options{})

	w := worker.New(c, fork_queue, worker.Options{})

	w.RegisterWorkflow(ForkWorkflow)
	w.RegisterWorkflow(BranchWorkflow)

	if err := w.Run(worker.InterruptCh()); err != nil {
		panic(err)
	}
}
