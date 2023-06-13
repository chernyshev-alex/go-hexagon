package workflow

import (
	"context"
	"fmt"
	"time"

	"log"

	"github.com/chernyshev-alex/go-hexagon/internal/domain/ports"
	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

const (
	SMS_QueueName = "send-sms-activity"
)

type sendToAuthorWorkflow struct {
	client      client.Client
	to, content string
}

func NewSendToAuthorWorkflow() *sendToAuthorWorkflow {
	return &sendToAuthorWorkflow{}
}

var _ ports.SmsSender = (*sendToAuthorWorkflow)(nil)

func (*sendToAuthorWorkflow) SendToAuthor(to string, content string) error {
	c, err := client.NewLazyClient(client.Options{HostPort: client.DefaultHostPort})
	if err != nil {
		panic(err.Error())
	}

	s := sendToAuthorWorkflow{client: c, to: to, content: content}
	defer s.client.Close()

	if err := s.startSendSmsWorkflow(); err != nil {
		return fmt.Errorf("SendToAuthor : failed to start workflow: %w", err)
	}
	if err := s.runSmsWorker(); err != nil {
		return fmt.Errorf("SendToAuthor : unable to run worker %w", err)
	}
	return nil
}

func (s sendToAuthorWorkflow) startSendSmsWorkflow() error {
	wOpts := client.StartWorkflowOptions{
		ID:        fmt.Sprintf("send-sms-%s-%s", s.to, uuid.NewString()),
		TaskQueue: SMS_QueueName,
	}
	wr, err := s.client.ExecuteWorkflow(context.Background(), wOpts, s.notifySmsWorkflow)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
		return err
	}
	log.Println("Started workflow", "WorkflowID", wr.GetID(), "RunID", wr.GetRunID())
	return nil
}

func (s sendToAuthorWorkflow) runSmsWorker() error {
	w := worker.New(s.client, SMS_QueueName, worker.Options{})
	w.RegisterWorkflow(s.notifySmsWorkflow)
	w.RegisterActivity(s.sendSmsActivity)

	err := w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Worker failed to start", err)
		return err
	}
	return nil
}

func (s sendToAuthorWorkflow) notifySmsWorkflow(ctx workflow.Context) error {
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

	if err := workflow.ExecuteActivity(ctx, s.sendSmsActivity).Get(ctx, nil); err != nil {
		logger.Error("Activity failed", "Error", err)
		return err
	}
	logger.Info("Completed")
	return nil
}

func (s sendToAuthorWorkflow) sendSmsActivity(ctx context.Context) error {
	// contact with external sms service
	log.Printf("SendSmsActivity %s\n", s.to)
	return nil
}
