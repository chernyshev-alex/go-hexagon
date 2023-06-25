package flows

import (
	"context"
	"fmt"
	"log"
	"time"

	models "github.com/chernyshev-alex/go-hexagon/internal/domain/model"
	"github.com/chernyshev-alex/go-hexagon/internal/domain/ports"
	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

type PublishProviderWorkflow struct{}
type ArticleCreatedWorklow struct{}
type ArticleRetrieverWorklow struct{}

var _ ports.PublishProvider = (*PublishProviderWorkflow)(nil)

func (p PublishProviderWorkflow) PublishArticleCreated(article *models.Article) error {
	return ArticleCreatedWorklow{}.Run(article)
}

func (p ArticleCreatedWorklow) Run(article *models.Article) error {
	const name = "article-creation"

	c, _ := client.Dial(client.Options{HostPort: client.DefaultHostPort})
	defer c.Close()

	opts := client.StartWorkflowOptions{
		ID:        fmt.Sprintf("%s-%s-%s", name, article.Id, uuid.NewString()),
		TaskQueue: name + "-queue",
	}
	if _, err := c.ExecuteWorkflow(context.Background(), opts, p.publishCreationOfWorkflow, article); err != nil {
		return err
	}

	w := worker.New(c, fork_queue, worker.Options{})
	w.RegisterWorkflow(p.publishCreationOfWorkflow)
	w.RegisterActivity(sendMessageForCreatedActivity)
	w.RegisterActivity(publishActivity)
	w.RegisterActivity(notifyAboutCreationOfActivity)

	go func() {
		if err := w.Run(worker.InterruptCh()); err != nil {
			log.Println(err)
		}
	}()
	return nil
}

func (p PublishProviderWorkflow) PublishArticleRetrieved(article *models.Article) error {
	return ArticleRetrieverWorklow{}.Run(article)
}

func (p ArticleRetrieverWorklow) Run(article *models.Article) error {
	const name = "article-retrieve"

	c, _ := client.Dial(client.Options{HostPort: client.DefaultHostPort})
	defer c.Close()

	opts := client.StartWorkflowOptions{
		ID:        fmt.Sprintf("%s-%s-%s", name, article.Id, uuid.NewString()),
		TaskQueue: name + "-queue",
	}
	if _, err := c.ExecuteWorkflow(context.Background(), opts, p.publishRetrievalOfWorkflow, article); err != nil {
		return err
	}

	w := worker.New(c, fork_queue, worker.Options{})
	w.RegisterWorkflow(p.publishRetrievalOfWorkflow)
	w.RegisterActivity(sendMessageForRetrieval)

	go func() {
		if err := w.Run(worker.InterruptCh()); err != nil {
			log.Println(err)
		}
	}()
	return nil
}

func (p ArticleCreatedWorklow) publishCreationOfWorkflow(ctx workflow.Context, article *models.Article) error {
	opts := workflow.ActivityOptions{StartToCloseTimeout: 5 * time.Second}
	logger := workflow.GetLogger(ctx)

	ctx1 := workflow.WithActivityOptions(ctx, opts)
	if err := workflow.ExecuteActivity(ctx1, sendMessageForCreatedActivity, article).Get(ctx1, nil); err != nil {
		logger.Error(err.Error())
	}
	if err := workflow.ExecuteActivity(ctx1, publishActivity).Get(ctx1, nil); err != nil {
		logger.Error(err.Error())
	}
	if err := workflow.ExecuteActivity(ctx1, notifyAboutCreationOfActivity).Get(ctx1, nil); err != nil {
		logger.Error(err.Error())
	}
	return nil
}

func (p ArticleRetrieverWorklow) publishRetrievalOfWorkflow(ctx workflow.Context, article *models.Article) error {
	opts := workflow.ActivityOptions{StartToCloseTimeout: 5 * time.Second}
	logger := workflow.GetLogger(ctx)

	ctx1 := workflow.WithActivityOptions(ctx, opts)
	if err := workflow.ExecuteActivity(ctx1, sendMessageForRetrieval, article).Get(ctx1, nil); err != nil {
		logger.Error(err.Error())
	}
	return nil
}

func sendMessageForCreatedActivity(ctx context.Context, article *models.Article) error {
	log.Println("sendMessageForCreatedActivity")
	return nil
}

func publishActivity(ctx context.Context, address string) error {
	log.Println("publishActivity", address)
	return nil
}

func notifyAboutCreationOfActivity(ctx context.Context, address string) error {
	log.Println("notifyAboutCreationOfActivity : ", address)
	return nil
}

func sendMessageForRetrieval(ctx context.Context, article *models.Article) error {
	log.Println("sendMessageForRetrieval")
	return nil
}
