package api

import (
	"context"

	"github.com/chernyshev-alex/go-hexagon/internal/domain/ports"
)

//go:generate go run github.com/vektra/mockery/v2@v2.32.4 --name ArticleFacade
type ArticleFacade interface {
	Get(ctx context.Context, articleId string) (ArticleResponse, error)
	Create(context.Context, ArticleRequest) (ArticleResponse, error)
	SearchBy(ctx context.Context, what, value string) ([]ArticleResponse, error)
}

type articleFacade struct {
	service ports.ArticleService
}

func NewArticleFacade(svc ports.ArticleService) ArticleFacade {
	return articleFacade{service: svc}
}

func (f articleFacade) Create(ctx context.Context, rq ArticleRequest) (ArticleResponse, error) {
	if article, err := f.service.Create(ctx, rq.AuthorId, rq.Title, rq.Content); err != nil {
		return ArticleResponse{}, err
	} else {
		return NewArticleResponse(&article), nil
	}
}

func (articleFacade) Get(ctx context.Context, articleId string) (ArticleResponse, error) {
	panic("unimplemented")
}

func (articleFacade) SearchBy(ctx context.Context, what string, value string) ([]ArticleResponse, error) {
	panic("unimplemented")
}
