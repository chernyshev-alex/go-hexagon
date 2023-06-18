package workflow

import (
	"log"
	"testing"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func TestForkJoinWorkflow(t *testing.T) {

	dag := DagDefinition{
		"fork-1": {
			"branch-1-1": "action-1-1",
			"branch-1-2": "action-1-2",
			"branch-1-3": "action-1-3",
			"branch-1-4": "action-1-4",
		},
		"fork-2": {
			"branch-2-1": "action-2-1",
			"branch-2-2": "action-2-2",
			"branch-2-3": "action-2-3",
		},
		"fork-3": {
			"branch-3-1": "action-3-1",
			"branch-3-2": "action-3-1",
		},
		"fork-4": {
			"branch-4-1": "action-4-1",
			"branch-4-2": "action-4-2",
			"branch-4-3": "action-4-3",
			"branch-4-4": "action-4-4",
		},
	}

	w := NewForkJoinWorkflow("INGEST_W22K", dag)
	if err := w.RunEtlWorkFlow(); err != nil {
		t.Fatalf("failed %s", err)
	}
}

// UI  http://localhost:8233/
// /usr/local/go/bin/go test -timeout 600s -run ^TestRunForkJoinWorker$ github.com/chernyshev-alex/go-hexagon/internal/adapters/workflow

func TestRunForkJoinWorker(t *testing.T) {
	cli, _ := client.Dial(client.Options{})
	defer cli.Close()

	workFlow := NewForkJoinWorkflow("INGEST_W22K", nil)
	// register worker's activities
	instance := worker.New(cli, workFlow.queueName(), worker.Options{})
	instance.RegisterWorkflow(workFlow.ParentWorkflow)
	instance.RegisterWorkflow(workFlow.TaskChildWorkflow)
	instance.RegisterActivity(workFlow.RunSQL)

	if err := instance.Run(worker.InterruptCh()); err != nil {
		log.Fatalln("Worker failed", err)
	}
	log.Println("finished OK")
}

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

	err := w.Run(worker.InterruptCh())
	if err != nil {
		panic(err)
	}
}
