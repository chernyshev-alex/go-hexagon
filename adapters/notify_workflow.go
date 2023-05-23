package adapters

import (
	"context"
	"fmt"
	"time"

	"log"

	"github.com/chernyshev-alex/go-hexagon/domain/ports"
	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

type sendSmsWorkflow struct {
	c client.Client
}

var _ ports.SmsSender = (*sendSmsWorkflow)(nil)

func CreateSender() *sendSmsWorkflow {
	c, err := client.NewLazyClient(client.Options{HostPort: client.DefaultHostPort})
	if err != nil {
		panic(err.Error())
	}
	return &sendSmsWorkflow{c: c}
}

func (s sendSmsWorkflow) CloseSender() error {
	s.c.Close()
	return nil
}

// TODO  imp different workflows

func (s sendSmsWorkflow) SendToAuthor(to string, content string) error {
	// TODO create specific workflow here
	if err := s.startSendSmsWorkflow(to); err != nil {
		return fmt.Errorf("failed to start SendSms workflow: %w", err)
	}
	if err := s.RunSmsWorker(); err != nil {
		log.Fatalln("Unable to run worker", err)
		return err
	}
	return nil
}

const (
	taskQueueName = "send-sms-activity"
)

func (s sendSmsWorkflow) startSendSmsWorkflow(sendTo string) error {
	wOpts := client.StartWorkflowOptions{
		ID:        fmt.Sprintf("send-sms-%s-%s", sendTo, uuid.NewString()),
		TaskQueue: taskQueueName,
	}
	wr, err := s.c.ExecuteWorkflow(context.Background(), wOpts, s.notifySmsWorkflow, sendTo)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
		return err
	}
	log.Println("Started workflow", "WorkflowID", wr.GetID(), "RunID", wr.GetRunID())
	return nil
}

func (s sendSmsWorkflow) notifySmsWorkflow(ctx workflow.Context, sendTo string) error {
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
	if err := workflow.ExecuteActivity(ctx, s.sendSmsActivity, sendTo).Get(ctx, nil); err != nil {
		logger.Error("Activity failed", "Error", err)
		return err
	}
	logger.Info("Completed")
	return nil
}

func (s sendSmsWorkflow) RunSmsWorker() error {
	w := worker.New(s.c, taskQueueName, worker.Options{})
	w.RegisterWorkflow(s.notifySmsWorkflow)
	w.RegisterActivity(s.sendSmsActivity)

	err := w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Worker failed to start", err)
		return err
	}
	return nil
}

func (s sendSmsWorkflow) sendSmsActivity(ctx context.Context, name string) error {
	fmt.Printf("SendSmsActivity %s", name)
	return nil
}
