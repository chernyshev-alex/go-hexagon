package adapters

import (
	"context"
	"time"

	"log"

	"github.com/chernyshev-alex/go-hexagon/domain/ports"
	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

type temporalSmsSender struct {
	c client.Client
}

var _ ports.SmsSender = (*temporalSmsSender)(nil)

func NewSmsSender() temporalSmsSender {
	if c, err := client.NewLazyClient(client.Options{
		HostPort: client.DefaultHostPort}); err != nil {
		panic(err.Error())
	} else {
		return temporalSmsSender{c: c}
	}
}

func (s temporalSmsSender) SendToAuthor(to string, content string) error {
	defer s.c.Close()
	s.StartSendSmsWorkflow()
	s.RunSmsWorker()
	return nil
}

func (s temporalSmsSender) StartSendSmsWorkflow() {
	workflowOptions := client.StartWorkflowOptions{
		ID:        "send-sms_activity_" + uuid.NewString(),
		TaskQueue: "send-sms-activity",
	}
	if wr, err := s.c.ExecuteWorkflow(context.Background(), workflowOptions, s.RetryNotifySmsWorkflow); err != nil {
		log.Fatalln("Unable to execute workflow", err)
	} else {
		log.Println("Started workflow", "WorkflowID", wr.GetID(), "RunID", wr.GetRunID())
	}
}

func (s temporalSmsSender) RunSmsWorker() {
	w := worker.New(s.c, "sms-activity", worker.Options{})
	w.RegisterWorkflow(s.RetryNotifySmsWorkflow)
	w.RegisterActivity(s.SendSmsActivity)

	err := w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}

func (s temporalSmsSender) RetryNotifySmsWorkflow(ctx workflow.Context) error {
	opts := workflow.ActivityOptions{
		StartToCloseTimeout: 2 * time.Minute,
		HeartbeatTimeout:    10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    5,
		},
	}

	ctx = workflow.WithActivityOptions(ctx, opts)
	logger := workflow.GetLogger(ctx)
	if err := workflow.ExecuteActivity(ctx, s.SendSmsActivity, 0, 20, time.Second).Get(ctx, nil); err != nil {
		logger.Error("Activity failed", "Error", err)
		return err
	}
	logger.Info("Completed")
	return nil
}

func (s temporalSmsSender) SendSmsActivity(ctx context.Context, name string) error {
	return nil
}
