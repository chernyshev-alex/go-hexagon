package ports

import "github.com/chernyshev-alex/go-hexagon/domain/models"

type ArticleMessageSender interface {
	SendMessageForCreated(*models.Article) error
	SendMessageForRetrieved(*models.Article) error
}
