package notify

import (
	"time"

	"github.com/chernyshev-alex/go-hexagon/internal/domain/model"
	"github.com/chernyshev-alex/go-hexagon/internal/domain/ports"
)

type ArticlePublisherRetriever interface {
	PublishRetrievalOf(*model.Article) error
	PublishCreationOf(*model.Article) error
	PublishArticleRetrieved(*model.Article) error
}

type ArticlePublisherProvider interface {
	PublishArticleCreated(*model.Article) error
	PublishArticleRetrieved(*model.Article) error
	ValidateEligibilityForPublication(*model.Article) error
}

type articlePublisher struct {
	sender    ports.ArticleMessageSender
	publisher ArticlePublisherProvider
	retriever ArticlePublisherRetriever
}

func NewArticlePublisher(
	s ports.ArticleMessageSender,
	p ArticlePublisherProvider,
	r ArticlePublisherRetriever) *articlePublisher {
	return &articlePublisher{
		sender:    s,
		publisher: p,
		retriever: r,
	}
}

func (p articlePublisher) PublishCreationOf(article *model.Article) error {
	return p.publisher.PublishArticleCreated(article)
}

func (p articlePublisher) PublishRetrievalOf(article *model.Article) error {
	return p.retriever.PublishArticleRetrieved(article)
}

// ====

type articleRetrievedMessage struct {
	article *model.Article
	sentAt  time.Time
}

func NewArticleRetrievedMessage(article *model.Article) articleCreatedMessage {
	return articleCreatedMessage{
		article: article,
		sentAt:  time.Now(),
	}
}

type articleCreatedMessage struct {
	article *model.Article
	sentAt  time.Time
}

func NewArticleCreatedMessage(article *model.Article) articleRetrievedMessage {
	return articleRetrievedMessage{
		article: article,
		sentAt:  time.Now(),
	}
}

type MsgBrokerArticleSender struct {
	ports.ArticleMessageSender
}

func (mb MsgBrokerArticleSender) SendMessageForCreated(article *model.Article) error {
	_ = NewArticleCreatedMessage(article)
	panic("not implemented")
}

func (mb MsgBrokerArticleSender) SendMessageForRetrieved(article *model.Article) error {
	_ = NewArticleRetrievedMessage(article)
	panic("not implemented")
}

type articleMailModel struct {
	recipientId string
	subject     string
	content     string
}

type authorMailNotifier struct {
	ports.AuthorNotifier
}

func NewArticleMailModel(a *model.Article) articleMailModel {
	return articleMailModel{
		recipientId: a.Author.Email,
		subject:     a.Title,
		content:     a.Content,
	}
}

func NewAuthorMailNotifier() authorMailNotifier {
	return authorMailNotifier{}
}

func (n authorMailNotifier) NotifyAboutCreationOf(article *model.Article) error {
	_ = NewArticleMailModel(article)
	panic("not implemented")
}

type articleSmsModel struct {
	recipientId string
	text        string
}

func NewArticleSmsModel(article *model.Article) articleSmsModel {
	return articleSmsModel{
		recipientId: article.Author.PhoneNumber,
		text:        article.Title,
	}
}

type authorSmsNotifier struct {
	ports.AuthorNotifier
}

func NewAuthorSmsNotifier() authorSmsNotifier {
	return authorSmsNotifier{}
}

func (n authorSmsNotifier) NotifyAboutCreationOf(article *model.Article) error {
	_ = NewArticleMailModel(article)
	panic("not immplemented")
}
