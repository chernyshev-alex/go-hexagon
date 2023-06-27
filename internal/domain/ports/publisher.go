package ports

import (
	models "github.com/chernyshev-alex/go-hexagon/internal/domain/model"
)

type ArticleMessageSender interface {
	SendMessageForCreated(*models.Article) error

	SendMessageForRetrieved(*models.Article) error
}

type SmsSender interface {
	SendToAuthor(to string, content string) error
}

type SocialMediaPublisher interface {
	Publish(*models.Article) error
}

type AuthorNotifier interface {
	NotifyAboutCreationOf(*models.Article) error
}

type ArticlePublisherRetriever interface {
	PublishRetrievalOf(*models.Article) error
	PublishCreationOf(*models.Article) error
}

type PublishProvider interface {
	PublishArticleCreated(*models.Article) error
	PublishArticleRetrieved(*models.Article) error
}

type ArticlePublisher struct {
	ArticlePublisherRetriever
	sender         ArticleMessageSender
	publisProvider PublishProvider
}

func NewArticlePublisher(ms ArticleMessageSender, pp PublishProvider) *ArticlePublisher {
	return &ArticlePublisher{
		sender:         ms,
		publisProvider: pp,
	}
}

func (p ArticlePublisher) PublishCreationOf(article *models.Article) error {
	return p.publisProvider.PublishArticleCreated(article)
}

func (p ArticlePublisher) PublishRetrievalOf(article *models.Article) error {
	return p.publisProvider.PublishArticleRetrieved(article)
}
