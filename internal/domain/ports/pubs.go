package ports

import "github.com/chernyshev-alex/go-hexagon/internal/domain/model"

type ArticleMessageSender interface {
	SendMessageForCreated(*model.Article) error
	SendMessageForRetrieved(*model.Article) error
}

type AuthorNotifier interface {
	NotifyAboutCreationOf(*model.Article) error
}

type SocialMediaPublisher interface {
	Publish(*model.Article) error
}
type SmsSender interface {
	SendToAuthor(to string, content string) error
}
