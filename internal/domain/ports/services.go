package ports

import (
	"strconv"

	"github.com/chernyshev-alex/go-hexagon/internal/domain/model"
)

type ArticleService struct {
	articleRepository ArticleRepository
	authorRepository  AuthorRepository
	publisher         ArticlePublisher
}

func NewArticleService(ar ArticleRepository, ur AuthorRepository, ap ArticlePublisher) *ArticleService {
	return &ArticleService{
		articleRepository: ar,
		authorRepository:  ur,
		publisher:         ap,
	}
}

func (s ArticleService) Create(authorId, title, content string) (string, error) {
	const EmptyStr = ""
	if author, err := s.authorRepository.Get(authorId); err != nil {
		return EmptyStr, err
	} else if article, err := s.articleRepository.Save(author, title, content); err != nil {
		return EmptyStr, err
	} else if err := article.ValidateEligibilityForPublication(); err != nil {
		return EmptyStr, err
	} else {
		s.publisher.PublishCreationOf(article)
		return article.Id, nil
	}
}

func (s ArticleService) Get(articleId string) (article *model.Article, err error) {
	dbId, _ := strconv.Atoi(articleId)
	if article, err = s.articleRepository.Get(uint(dbId)); err != nil {
		return nil, err
	}
	err = s.publisher.PublishRetrievalOf(article)
	return article, err
}
