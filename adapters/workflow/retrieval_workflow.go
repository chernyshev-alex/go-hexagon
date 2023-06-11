package workflow

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
)

type ForkJoinWorkflow struct {
	wkName string
	config map[string]interface{}
}

func NewForkJoinWorkflow(name string, config map[string]interface{}) *ForkJoinWorkflow {
	return &ForkJoinWorkflow{wkName: name, config: config}
}

func (p ForkJoinWorkflow) queueName() string { return p.wkName + "_queue" }

func (p ForkJoinWorkflow) RunEtlWorkFlow() error {
	cli, _ := client.Dial(client.Options{})
	defer cli.Close()

	opts := client.StartWorkflowOptions{
		ID:        fmt.Sprintf("etl-%s-%s", p.wkName, uuid.NewString()),
		TaskQueue: p.queueName(),
	}

	var results []string

	workflowRun, _ := cli.ExecuteWorkflow(context.Background(), opts, p.ParentWorkflow)
	if err := workflowRun.Get(context.Background(), &results); err != nil {
		return err
	}
	return nil
}

func (p ForkJoinWorkflow) ParentWorkflow(ctx workflow.Context) ([]string, error) {
	futures := []workflow.Future{}
	for i := 0; i < 5; i++ {
		childCtx := workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{})
		future := workflow.ExecuteChildWorkflow(childCtx, p.TaskChildWorkflow, fmt.Sprintf("Task-%d", i), p.config)
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

func (p ForkJoinWorkflow) TaskChildWorkflow(ctx workflow.Context, taskID string, config map[string]interface{}) (string, error) {
	opts := workflow.ActivityOptions{StartToCloseTimeout: 10 * time.Second}
	ctx1 := workflow.WithActivityOptions(ctx, opts)
	logger := workflow.GetLogger(ctx1)

	var result string
	var sql = "insert into .." //  get SQL to execute
	if err := workflow.ExecuteActivity(ctx1, p.RunSQL, taskID, sql).Get(ctx1, &result); err != nil {
		logger.Error("failed", err)
		return "", err
	}
	return taskID + ":" + result, nil
}

var task3_count = 3

func (p *ForkJoinWorkflow) RunSQL(ctx context.Context, taskID string, sqlStatement string) (string, error) {
	// connect to db and run SQL statement
	if taskID == "Task-3" && task3_count > 0 {
		task3_count--
		return "", fmt.Errorf("emulate error in task-3")
	}
	time.Sleep(5 * time.Second)
	return "OK", nil
}
