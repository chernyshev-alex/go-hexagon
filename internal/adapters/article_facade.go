package adapters

import (
	"github.com/chernyshev-alex/go-hexagon/internal/domain/ports"
)

type ArticleFacade interface {
	Get(articleId string) (*ArticleResponse, error)
	Create(rq *ArticleRequest) (ArticleIdResponse, error)
	List(pageId string) ([]*ArticleResponse, error)
	SearchBy(what, value string) ([]*ArticleResponse, error)
}

var _ ArticleFacade = (*ArticleFacadeImpl)(nil)

type ArticleFacadeImpl struct {
	articleService ports.ArticleService
}

func NewArticleFacade(service ports.ArticleService) ArticleFacadeImpl {
	return ArticleFacadeImpl{
		articleService: service,
	}
}

func (f ArticleFacadeImpl) Get(articleId string) (*ArticleResponse, error) {
	article, err := f.articleService.Get(articleId)
	if err != nil {
		return nil, err
	}
	return NewArticleResponseWith().Article(article), nil
}

func (f ArticleFacadeImpl) Create(rq *ArticleRequest) (ArticleIdResponse, error) {
	articleId, err := f.articleService.Create(rq.AuthorId, rq.Title, rq.Content)
	if err != nil {
		return ArticleIdResponse{}, err
	}
	return NewArticleIdResponse(articleId), nil
}

func (f ArticleFacadeImpl) List(pageId string) ([]*ArticleResponse, error) {
	return []*ArticleResponse{}, nil
}

func (f ArticleFacadeImpl) SearchBy(what, value string) ([]*ArticleResponse, error) {
	return []*ArticleResponse{}, nil
}
