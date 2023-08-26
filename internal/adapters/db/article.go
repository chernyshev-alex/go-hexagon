package adapters

import (
	"context"
	"strconv"

	"github.com/chernyshev-alex/go-hexagon/internal/domain/model"
	"github.com/chernyshev-alex/go-hexagon/internal/domain/ports"
	"gorm.io/gorm"
)

type articleTable struct {
	gorm.Model
	Title      string
	Content    string
	AuthorName string
	AuthorId   string
}

func (s *articleTable) TableName() string {
	return "article"
}

func (s *articleTable) ToDomain() model.Article {
	authorId, _ := strconv.Atoi(s.AuthorId)
	return model.Article{
		ID:      int64(s.ID),
		Title:   s.Title,
		Content: s.Content,
		Author: model.Author{
			ID:   int64(authorId),
			Name: s.AuthorName,
		},
	}
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&articleTable{})
}

type DbArticleRepository struct {
	db *gorm.DB
}

func NewDbArticleRepository(db *gorm.DB) ports.ArticleRepository {
	return DbArticleRepository{db: db}
}

func (r DbArticleRepository) GetByID(ctx context.Context, id int64) (model.Article, error) {
	var t_article articleTable

	if err := r.db.First(&t_article, id).Error; err != nil {
		return model.Article{}, err
	}
	return t_article.ToDomain(), nil
}

func (DbArticleRepository) DeleteById(context.Context, int64) error {
	panic("unimplemented")
}

func (DbArticleRepository) Fetch(ctx context.Context, cursor string, num int64) (xs []model.Article, nextCursor string, err error) {
	panic("unimplemented")
}

func (DbArticleRepository) GetByTitle(context.Context, string) (model.Article, error) {
	panic("unimplemented")
}

func (DbArticleRepository) Store(context.Context, *model.Article) error {
	panic("unimplemented")
}

func (DbArticleRepository) Update(context.Context, *model.Article) error {
	panic("unimplemented")
}

// func (r *DbArticleRepository) Get(id uint) (*model.Article, error) {
// 	var db_article articleStorage
// 	if err := r.db.First(&db_article, id).Error; err != nil {
// 		return nil, err
// 	}
// 	return db_article.ToDomain(), nil
// }

// func (r *DbArticleRepository) Save(author *model.Author, title string, content string) (*model.Article, error) {
// 	db_model := NewArticleStorage(author, title, DbArticleModelWithContent(content))
// 	if err := r.db.Create(&db_model).Error; err != nil {
// 		return nil, err
// 	}
// 	return db_model.ToDomain(), nil
// }
