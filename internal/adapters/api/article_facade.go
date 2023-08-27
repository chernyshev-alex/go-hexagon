package api

import (
	"context"
	"strconv"

	"github.com/chernyshev-alex/go-hexagon/internal/domain/ports"
)

//go:generate go run github.com/vektra/mockery/v2@v2.33.0 --name ArticleFacade
type ArticleFacade interface {
	Create(context.Context, ArticleRequest) (ArticleResponse, error)
	Get(ctx context.Context, articleId string) (ArticleResponse, error)
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

func (f articleFacade) Get(ctx context.Context, articleId string) (ArticleResponse, error) {
	id := mustBeInt(articleId)
	if article, err := f.service.GetByID(ctx, int64(id)); err != nil {
		return ArticleResponse{}, err
	} else {
		return NewArticleResponse(&article), nil
	}
}

func (articleFacade) SearchBy(ctx context.Context, what string, value string) ([]ArticleResponse, error) {
	panic("unimplemented")
}

func mustBeInt(articleId string) int {
	if id, err := strconv.Atoi(articleId); err != nil {
		panic(err)
	} else {
		return id
	}
}
