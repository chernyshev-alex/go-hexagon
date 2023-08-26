package api

import "github.com/chernyshev-alex/go-hexagon/internal/domain/model"

type ArticleRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	AuthorId int64  `json:"authorid"`
}

type ArticleResponse struct {
	ID         int64  `json:"id"`
	Title      string `json:"title"`
	AuthorName string `json:"authorName"`
}

func NewArticleResponse(a *model.Article) ArticleResponse {
	return ArticleResponse{
		ID:         a.ID,
		Title:      a.Title,
		AuthorName: a.Author.Name,
	}
}
