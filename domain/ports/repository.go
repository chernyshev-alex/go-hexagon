package ports

import m "github.com/chernyshev-alex/go-hexagon/domain/models"

type ArticleRepository interface {
	Save(author *m.Author, title string, content string) (*m.Article, error)
	Get(uint) (*m.Article, error)
}

type AuthorRepository interface {
	Get(authorId string) (*m.Author, error)
}
