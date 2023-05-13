package models

import "strconv"

type Author struct {
	Id, Name string
}

func NewAuthor(id, name string) *Author {
	return &Author{
		Id:   id,
		Name: name,
	}
}

type Article struct {
	Id         string
	Title      string
	Content    string
	AuthorName string
}

type ArticleOpt func(*Article)

func NewArticleOpt(opts ...ArticleOpt) *Article {
	a := &Article{}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

func ArticleWithId(id uint) ArticleOpt {
	return func(a *Article) {
		a.Id = strconv.Itoa(int(id))
	}
}

func ArticleWithAuthor(authorName string) ArticleOpt {
	return func(a *Article) {
		a.AuthorName = authorName
	}
}

func ArticleWithTitle(title string) ArticleOpt {
	return func(a *Article) {
		a.Title = title
	}
}

func ArticleWithContent(content string) ArticleOpt {
	return func(a *Article) {
		a.Content = content
	}
}

func (a *Article) ValidateEligibilityForPublication() error {
	return nil
}
