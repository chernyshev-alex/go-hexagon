package adapters

import (
	"fmt"

	"github.com/chernyshev-alex/go-hexagon/domain/models"
	"github.com/chernyshev-alex/go-hexagon/domain/ports"
)

const TWEET = "Check out the new article >>%s<< by %s"

type ArticleTwitterModel struct {
	twitterAccount string
	tweet          string
}

func NewArticleTwitterModel(twitterAccount string, tweet string) *ArticleTwitterModel {
	return &ArticleTwitterModel{
		twitterAccount: twitterAccount,
		tweet:          tweet,
	}
}

func NewArticleTwitterModelFromArticle(a *models.Article) *ArticleTwitterModel {
	return NewArticleTwitterModel(a.AuthorName, fmt.Sprintf(TWEET, a.Title, a.AuthorName))
}

func (m ArticleTwitterModel) ToString() string {
	return m.tweet
}

type TwitterClient struct{}

func (tc TwitterClient) tweet(a *ArticleTwitterModel) error {
	return nil
}

type TwitterArticlePublisher struct {
	client *TwitterClient
}

var _ ports.SocialMediaPublisher = (*TwitterArticlePublisher)(nil)

func NewTwitterArticlePublisher(c *TwitterClient) *TwitterArticlePublisher {
	return &TwitterArticlePublisher{
		client: c,
	}
}

func (p *TwitterArticlePublisher) Publish(a *models.Article) error {
	articleTweet := NewArticleTwitterModelFromArticle(a)
	p.client.tweet(articleTweet)
	return nil
}
