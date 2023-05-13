package ports

import m "github.com/chernyshev-alex/go-hexagon/domain/models"

type SocialMediaPublisher interface {
	Publish(*m.Article) error
}

type AuthorNotifier interface {
	NotifyAboutCreationOf(*m.Article) error
}

type ArticlePublisher struct {
	sender     ArticleMessageSender
	publishers []SocialMediaPublisher
	notifiers  []AuthorNotifier
}

func NewArticlePublisher(ms ArticleMessageSender, pubs []SocialMediaPublisher, nfs []AuthorNotifier) *ArticlePublisher {
	return &ArticlePublisher{
		sender:     ms,
		publishers: pubs,
		notifiers:  nfs,
	}
}

func (p ArticlePublisher) PublishCreationOf(article *m.Article) (err error) {
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

func (p ArticlePublisher) publishRetrievalOf(article *m.Article) error {
	return p.sender.SendMessageForRetrieved(article)
}
