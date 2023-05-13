package ports

import m "github.com/chernyshev-alex/go-hexagon/domain/models"

type ArticleMessageSender interface {
	SendMessageForCreated(*m.Article) error
	SendMessageForRetrieved(*m.Article) error
}
