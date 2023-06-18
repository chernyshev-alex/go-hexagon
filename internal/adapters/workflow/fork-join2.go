package workflow

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

var fork_queue = "my-task-queue"

var _ = rand.New(rand.NewSource(32))

func BranchWorkflow(ctx workflow.Context, message string) (string, error) {
	workflow.Sleep(ctx, time.Duration(rand.Intn(20))*time.Second)
	return message, nil
}

func ForkWorkflow(ctx workflow.Context, pipeline map[string][]string, order []string) ([]string, error) {
	results := make([]string, 0)

	for _, fork := range order {
		branches := pipeline[fork]
		doneCh := workflow.NewChannel(ctx)

		for _, branch := range branches {
			branch := branch

			workflow.Go(ctx, func(ctx workflow.Context) {
				ctx1 := workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
					WorkflowID:  branch,
					RetryPolicy: &temporal.RetryPolicy{},
					// CronSchedule: "* * * * *",
				})

				var result string
				if err := workflow.ExecuteChildWorkflow(ctx1, BranchWorkflow, branch).Get(ctx1, &result); err != nil {
					fmt.Printf("Error executing branch %s: %v\n", branch, err)
				} else {
					results = append(results, result)
				}

				doneCh.Send(ctx, nil)
			})
		}
		for range branches {
			doneCh.Receive(ctx, nil)
		}
	}
	return results, nil
}

func RunWorflow(pipeline map[string][]string, order []string) {
	c, _ := client.Dial(client.Options{})

	ctx := context.Background()
	MustWorkflow(func() (client.WorkflowRun, error) {
		options := client.StartWorkflowOptions{TaskQueue: fork_queue}
		return c.ExecuteWorkflow(ctx, options, ForkWorkflow, pipeline, order)
	}).Get(ctx, nil)
}

func MustWorkflow(wf func() (client.WorkflowRun, error)) client.WorkflowRun {
	if f, err := wf(); err != nil {
		panic(err)
	} else {
		return f
	}
}
