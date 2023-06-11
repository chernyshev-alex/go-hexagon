package ports

import (
	models "github.com/chernyshev-alex/go-hexagon/domain/model"
)

type ArticleRepository interface {
	Save(author *models.Author, title string, content string) (*models.Article, error)

	Get(uint) (*models.Article, error)
}

type AuthorRepository interface {
	Get(authorId string) (*models.Author, error)
}
