package adapters

import (
	"errors"

	"github.com/chernyshev-alex/go-hexagon/internal/domain/model"
)

type Verror struct {
	error
}

func MakeVerror(msg string) error {
	return &Verror{errors.New(msg)}
}

type ArticleRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	AuthorId string `json:"authorid"`
}

func NewArticleRequest(authorId, title, content string) *ArticleRequest {
	return &ArticleRequest{
		Title:    title,
		Content:  content,
		AuthorId: authorId,
	}
}

func (r ArticleRequest) ValidateArticleRequest(rq *ArticleRequest) error {
	if rq.AuthorId == "" {
		return MakeVerror("authorId is required")
	}
	return nil
}

type ArticleIdResponse struct {
	Id string
}

func NewArticleIdResponse(id string) ArticleIdResponse {
	return ArticleIdResponse{
		Id: id,
	}
}

type ArticleResponse struct {
	Id         string `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	AuthorName string `json:"authorname"`
}

func NewArticleResponse(id, title, content, authorName string) *ArticleResponse {
	return &ArticleResponse{
		Id:         id,
		Title:      title,
		Content:    content,
		AuthorName: authorName,
	}
}

func NewArticleResponseWith() *ArticleResponse {
	return &ArticleResponse{}
}

func (ar *ArticleResponse) Article(article *model.Article) *ArticleResponse {
	ar.Id = article.Id
	ar.AuthorName = article.AuthorName
	ar.Content = article.Content
	ar.Title = article.Title
	return ar
}
