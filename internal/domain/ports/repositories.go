package ports

import (
	"context"

	"github.com/chernyshev-alex/go-hexagon/internal/domain/model"
)

type AuthorRepository interface {
	GetByID(ctx context.Context, id int64) (model.Author, error)
}

type ArticleRepository interface {
	Fetch(ctx context.Context, cursor string, num int64) (xs []model.Article, nextCursor string, err error)
	GetByID(context.Context, int64) (model.Article, error)
	GetByTitle(context.Context, string) (model.Article, error)
	Update(context.Context, *model.Article) error
	Store(context.Context, *model.Article) error
	DeleteById(context.Context, int64) error
}
