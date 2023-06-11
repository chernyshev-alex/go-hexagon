package workflow

import (
	"log"
	"testing"

	models "github.com/chernyshev-alex/go-hexagon/domain/model"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func TestPublishWorkflow(t *testing.T) {
	w := new(PublishCreationOfWorkflow)

	var article models.Article
	err := w.PublishCreationOf(&article)
	if err != nil {
		t.Fatalf("failed %s", err)
	}
}

func TestPublishWorker(t *testing.T) {
	cli, _ := client.NewLazyClient(client.Options{HostPort: client.DefaultHostPort})

	wk := worker.New(cli, publish_article_queue, worker.Options{})
	defer cli.Close()

	w := new(PublishCreationOfWorkflow)
	wk.RegisterWorkflow(w.publishCreationOfWorkflow)

	wk.RegisterActivity(sendMessageForCreatedActivity)
	wk.RegisterActivity(getPublishListActivity)
	wk.RegisterActivity(publishActivity)
	wk.RegisterActivity(notifyAboutCreationOfActivity)

	if err := wk.Run(worker.InterruptCh()); err != nil {
		log.Fatalln("Worker failed", err)
	}
	log.Println("finished")
}
