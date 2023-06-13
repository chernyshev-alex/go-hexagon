package ports

import (
	"github.com/chernyshev-alex/go-hexagon/internal/domain/model"
)

type ArticleRepository interface {
	Save(author *model.Author, title string, content string) (*model.Article, error)
	Get(uint) (*model.Article, error)
}

type AuthorRepository interface {
	Get(authorId string) (*model.Author, error)
}
