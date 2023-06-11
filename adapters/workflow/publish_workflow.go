package workflow

import (
	"context"
	"fmt"
	"log"
	"time"

	models "github.com/chernyshev-alex/go-hexagon/domain/model"
	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
)

//  Temporal.io workflow implementation

const publish_article_queue = "publish_article"

type PublishCreationOfWorkflow struct {
}

func (p PublishCreationOfWorkflow) PublishCreationOf(article *models.Article) error {
	c, err := client.Dial(client.Options{HostPort: client.DefaultHostPort})
	if err != nil {
		panic(err.Error())
	}

	defer c.Close()
	opts := client.StartWorkflowOptions{
		ID:        fmt.Sprintf("pub-creation-%s-%s", article.Id, uuid.NewString()),
		TaskQueue: publish_article_queue,
	}
	wr, err := c.ExecuteWorkflow(context.Background(), opts, p.publishCreationOfWorkflow, article)
	if err != nil {
		return err
	}
	log.Println("Started workflow", "WorkflowID", wr.GetID(), "RunID", wr.GetRunID())
	return nil
}

type Addresses struct {
	Pubs   []string
	Notify []string
}

func (p PublishCreationOfWorkflow) publishCreationOfWorkflow(ctx workflow.Context, article *models.Article) error {
	opts := workflow.ActivityOptions{StartToCloseTimeout: 5 * time.Second}
	logger := workflow.GetLogger(ctx)

	ctx1 := workflow.WithActivityOptions(ctx, opts)
	err := workflow.ExecuteActivity(ctx1, sendMessageForCreatedActivity, article).Get(ctx1, nil)
	if err != nil {
		logger.Error(err.Error())
	}

	var addresses Addresses
	err = workflow.ExecuteActivity(ctx1, getPublishListActivity, article).Get(ctx1, &addresses)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	var futures []workflow.Future
	for _, pub := range addresses.Pubs {
		f := workflow.ExecuteActivity(ctx1, publishActivity, pub)
		futures = append(futures, f)
	}
	for _, notify := range addresses.Notify {
		f := workflow.ExecuteActivity(ctx1, notifyAboutCreationOfActivity, notify)
		futures = append(futures, f)
	}
	// wait all futures
	for _, future := range futures {
		_ = future.Get(ctx, nil)
	}
	logger.Info("Workflow completed")
	return nil
}

func sendMessageForCreatedActivity(ctx context.Context, article *models.Article) error {
	log.Println("sendMessageForCreatedActivity")
	return nil
}

func getPublishListActivity(ctx context.Context, article *models.Article) (Addresses, error) {
	return Addresses{
		Pubs:   []string{"pub1", "pub2", "pub3"},
		Notify: []string{"notify1", "notify2"},
	}, nil
}

func publishActivity(ctx context.Context, address string) error {
	log.Println("publishActivity", address)
	return nil
}

func notifyAboutCreationOfActivity(ctx context.Context, address string) error {
	log.Println("notifyAboutCreationOfActivity : ", address)
	return nil
}
