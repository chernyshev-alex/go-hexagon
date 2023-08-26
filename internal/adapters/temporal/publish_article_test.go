package flows

import (
	"testing"

	models "github.com/chernyshev-alex/go-hexagon/internal/domain/model"
)

func TestPublishWorkflow(t *testing.T) {
	w := new(PublishProviderWorkflow)

	var article models.Article
	err := w.PublishArticleCreated(&article)
	if err != nil {
		t.Fatalf("failed %s", err)
	}
}
