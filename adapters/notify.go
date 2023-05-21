package adapters

import (
	"fmt"

	"github.com/chernyshev-alex/go-hexagon/domain/models"
	"github.com/chernyshev-alex/go-hexagon/domain/ports"
)

const (
	MAIL_SUBJECT = "You have successfully published: >>%s<<"
	MAIL_CONTENT = "Check if everything is correct: >>%s<<"
)

type ArticleMailModel struct {
	recipientId string
	subject     string
	content     string
}

func NewArticleMailModel(recipientId, subject, content string) *ArticleMailModel {
	return &ArticleMailModel{
		recipientId: recipientId,
		subject:     subject,
		content:     content,
	}
}
func NewArticleMailModelFromArticle(a *models.Article) *ArticleMailModel {
	return NewArticleMailModel(a.AuthorName,
		fmt.Sprintf(MAIL_SUBJECT, a.Title),
		fmt.Sprintf(MAIL_CONTENT, a.Content))
}
func (m ArticleMailModel) ToString() string {
	return m.content
}

type ArticleSmsModel struct {
	recipientId string
	text        string
}

const SMS_CONTENT = "Please check your email. We have sent you publication details of the article: >>%s<<"

func NewArticleSmsModel(recipientId, text string) *ArticleSmsModel {
	return &ArticleSmsModel{
		recipientId: recipientId,
		text:        text,
	}
}

func NewArticleSmsModelFromArticle(a *models.Article) *ArticleSmsModel {
	return &ArticleSmsModel{a.AuthorName, fmt.Sprintf(SMS_CONTENT, a.Title)}
}
func (m ArticleSmsModel) ToString() string {
	return m.text
}

type AuthorMailNotifier struct {
	client interface{}
}

var _ ports.AuthorNotifier = (*AuthorMailNotifier)(nil)

func NewAuthorMailNotifier(c interface{}) *AuthorMailNotifier {
	return &AuthorMailNotifier{
		client: c,
	}
}

func (n AuthorMailNotifier) NotifyAboutCreationOf(a *models.Article) error {
	_ = NewArticleMailModelFromArticle(a)
	// TODO notify
	return nil
}

type AuthorSmsNotifier struct {
	client ports.SmsSender
}

func NewAuthorSmsNotifier(c ports.SmsSender) *AuthorSmsNotifier {
	return &AuthorSmsNotifier{
		client: c,
	}
}
func (n *AuthorSmsNotifier) NotifyAboutCreationOf(a *models.Article) error {
	_ = NewArticleSmsModelFromArticle(a)
	// TODO notify
	return nil
}
