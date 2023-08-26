package ports

import (
	"context"
	"time"

	"github.com/chernyshev-alex/go-hexagon/internal/domain/model"
)

type ArticleService interface {
	Create(ctx context.Context, authorId int64, title, content string) (model.Article, error)
	Fetch(ctx context.Context, cursor string, num int64) ([]model.Article, string, error)
	GetByID(ctx context.Context, id int64) (model.Article, error)
	Update(ctx context.Context, ar *model.Article) error
	GetByTitle(ctx context.Context, title string) (model.Article, error)
	Store(context.Context, *model.Article) error
	Delete(ctx context.Context, id int64) error
}

type articleService struct {
	repoArticle ArticleRepository
	repoAuthor  AuthorRepository
	publisher   SocialMediaPublisher
	ctxTimeout  time.Duration
}

func NewArticleService(
	ar ArticleRepository,
	ur AuthorRepository,
	ap SocialMediaPublisher,
	t time.Duration) ArticleService {

	return articleService{
		repoArticle: ar,
		repoAuthor:  ur,
		publisher:   ap,
		ctxTimeout:  t,
	}
}

func (s articleService) Create(ctx context.Context, authorId int64, title, content string) (model.Article, error) {
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()

	author, err := s.repoAuthor.GetByID(ctx, authorId)
	if err != nil {
		return model.Article{}, nil
	}

	article := model.Article{
		Title:   title,
		Content: content,
		Author:  author,
	}
	if err = s.repoArticle.Store(ctx, &article); err != nil {
		return model.Article{}, nil
	}
	return article, nil
}

func (s articleService) Fetch(ctx context.Context, cursor string, nm int64) ([]model.Article, string, error) {
	// 	g, gctx := errgroup.WithContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()

	return s.repoArticle.Fetch(ctx, cursor, nm)
}

func (s articleService) GetByID(ctx context.Context, articleId int64) (model.Article, error) {
	if article, err := s.repoArticle.GetByID(ctx, articleId); err != nil {
		return model.Article{}, err
	} else {
		err = s.publisher.Publish(&article)
		return article, err
	}
}

func (articleService) GetByTitle(ctx context.Context, title string) (model.Article, error) {
	panic("unimplemented")
}

func (articleService) Store(context.Context, *model.Article) error {
	panic("unimplemented")
}

func (articleService) Update(ctx context.Context, ar *model.Article) error {
	panic("unimplemented")
}

func (articleService) Delete(ctx context.Context, id int64) error {
	panic("unimplemented")
}
