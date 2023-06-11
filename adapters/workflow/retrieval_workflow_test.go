package workflow

import (
	"log"
	"testing"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func TestForkJoinWorkflow(t *testing.T) {
	w := NewForkJoinWorkflow("INGEST_W22K", map[string]interface{}{})
	if err := w.RunEtlWorkFlow(); err != nil {
		t.Fatalf("failed %s", err)
	}
}

// UI  http://localhost:8233/

func TestForkJoinWorker(t *testing.T) {
	cli, _ := client.Dial(client.Options{})
	defer cli.Close()

	workFlow := NewForkJoinWorkflow("INGEST_W22K", map[string]interface{}{})
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
