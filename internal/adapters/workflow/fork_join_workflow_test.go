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
			"branch-1": "action-1",
			"branch-2": "action-2",
			"branch-3": "action-3",
			"branch-4": "action-4",
		},
		"fork-2": {
			"branch-5": "action-5",
			"branch-6": "action-6",
			"branch-7": "action-7",
		},
	}

	w := NewForkJoinWorkflow("INGEST_W22K", dag)
	if err := w.RunEtlWorkFlow(); err != nil {
		t.Fatalf("failed %s", err)
	}
}

// UI  http://localhost:8233/

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
	log.Println("finished")
}
