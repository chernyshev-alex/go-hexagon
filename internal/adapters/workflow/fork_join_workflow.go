package workflow

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type DagDefinition = map[string]map[string]string

type ForkJoinWorkflow struct {
	name string
	dag  DagDefinition
}

func NewForkJoinWorkflow(name string, dag DagDefinition) *ForkJoinWorkflow {
	return &ForkJoinWorkflow{name: name, dag: dag}
}

func (p ForkJoinWorkflow) queueName() string { return "fork-queue" }

func (p ForkJoinWorkflow) RunEtlWorkFlow() error {
	cli, _ := client.Dial(client.Options{})
	defer cli.Close()

	var results = make([]string, 0)
	for k, v := range p.dag {
		opts := client.StartWorkflowOptions{
			ID:        fmt.Sprintf("%s-%s", k, uuid.NewString()),
			TaskQueue: p.queueName(),
		}
		ctx := context.Background()
		workflowRun, _ := cli.ExecuteWorkflow(ctx, opts, p.ParentWorkflow, v)
		if err := workflowRun.Get(ctx, &results); err != nil {
			fmt.Println("******  ", k, err)
			return err
		}
	}
	return nil
}

func (p ForkJoinWorkflow) ParentWorkflow(ctx workflow.Context, branch map[string]string) ([]string, error) {
	futures := []workflow.Future{}
	for k, v := range branch {
		childCtx := workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
			WorkflowID:  k,
			RetryPolicy: &temporal.RetryPolicy{},
			// CronSchedule: "* * * * *",
		})
		future := workflow.ExecuteChildWorkflow(childCtx, p.TaskChildWorkflow, fmt.Sprintf("Task-%s", k), v)
		futures = append(futures, future)
	}
	logger := workflow.GetLogger(ctx)
	results := make([]string, len(futures))
	for i, f := range futures {
		var result string
		if err := f.Get(ctx, &result); err != nil {
			logger.Error(err.Error())
		}
		results[i] = result
	}
	return results, nil
}

func (p ForkJoinWorkflow) TaskChildWorkflow(ctx workflow.Context, taskID string, action string) (string, error) {
	opts := workflow.ActivityOptions{StartToCloseTimeout: 30 * time.Second}
	ctx1 := workflow.WithActivityOptions(ctx, opts)
	logger := workflow.GetLogger(ctx1)

	var result string
	if err := workflow.ExecuteActivity(ctx1, p.RunSQL, taskID, action).Get(ctx1, &result); err != nil {
		logger.Error("failed", err)
		return "", err
	}
	return taskID + ":" + result, nil
}

var task3_count = 3
var _ = rand.New(rand.NewSource(32))

func (p *ForkJoinWorkflow) RunSQL(ctx context.Context, taskID string, action string) (string, error) {
	// connect to db and run SQL statement
	// if taskID == "Task-branch-3" && task3_count > 0 {
	// 	task3_count--
	// 	return "", fmt.Errorf("emulate error in task-3")
	// }

	for {
		//  call db service
		if rand.Intn(5) == 3 {
			return "", nil
		}

		activity.RecordHeartbeat(ctx)
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(5 * time.Second):
		}
	}
	// log.Println("run activity", taskID, action)
	// return "OK", nil
}
