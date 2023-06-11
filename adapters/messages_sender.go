package adapters

import (
	"fmt"
	"time"

	"github.com/chernyshev-alex/go-hexagon/domain/model"
	"github.com/chernyshev-alex/go-hexagon/ports"
)

type ArticleCreatedMessage struct {
	article *model.Article
	ts      time.Time
}

func NewArticleCreatedMsg(a *model.Article) *ArticleCreatedMessage {
	return &ArticleCreatedMessage{
		article: a,
		ts:      time.Now(),
	}
}

func (acm ArticleCreatedMessage) ToString() string {
	return fmt.Sprintf("<<<article %s created>>>", acm.article.Title)
}

type ArticleRetrievedMessage struct {
	article *model.Article
	ts      time.Time
}

func NewArticleRetrievedMsg(article *model.Article) *ArticleRetrievedMessage {
	return &ArticleRetrievedMessage{
		article: article,
		ts:      time.Now(),
	}
}

func (acm ArticleRetrievedMessage) ToString() string {
	return fmt.Sprintf("<<<article %s retrieved>>>", acm.article.Title)

}

type MbArticleMessageSender struct {
	mbClient interface{}
}

var _ ports.ArticleMessageSender = (*MbArticleMessageSender)(nil)

func NewMbArticleMessageSender(client interface{}) *MbArticleMessageSender {
	return &MbArticleMessageSender{mbClient: client}
}

func (mb MbArticleMessageSender) SendMessageForCreated(a *model.Article) error {
	_ = NewArticleCreatedMsg(a) // send msg here
	return nil
}

func (mb MbArticleMessageSender) SendMessageForRetrieved(a *model.Article) error {
	_ = NewArticleRetrievedMsg(a)
	return nil
}
