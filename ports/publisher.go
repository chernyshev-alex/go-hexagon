package ports

import (
	models "github.com/chernyshev-alex/go-hexagon/domain/model"
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

type ArticlePublisher struct {
	ArticlePublisherRetriever
	sender     ArticleMessageSender
	publishers []SocialMediaPublisher
	notifiers  []AuthorNotifier
}

func NewArticlePublisher(ms ArticleMessageSender,
	pubs []SocialMediaPublisher,
	nfs []AuthorNotifier) *ArticlePublisher {

	return &ArticlePublisher{
		sender:     ms,
		publishers: pubs,
		notifiers:  nfs,
	}
}

func (p ArticlePublisher) PublishCreationOf(article *models.Article) (err error) {
	if err = p.sender.SendMessageForCreated(article); err != nil {
		return err
	}
	for _, pub := range p.publishers {
		err = pub.Publish(article)
	}
	for _, ntf := range p.notifiers {
		err = ntf.NotifyAboutCreationOf(article)
	}
	return err
}

func (p ArticlePublisher) publishRetrievalOf(article *models.Article) error {
	return p.sender.SendMessageForRetrieved(article)
}
