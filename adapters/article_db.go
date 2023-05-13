package adapters

import (
	"github.com/chernyshev-alex/go-hexagon/domain/models"
	"github.com/chernyshev-alex/go-hexagon/domain/ports"
	"gorm.io/gorm"
)

type ArticleDatabaseModel struct {
	gorm.Model
	Title      string
	Content    string
	AuthorName string
}

func (a *ArticleDatabaseModel) TableName() string {
	return "article"
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&ArticleDatabaseModel{})
}

type ArticleDatabaseModelOpt func(*ArticleDatabaseModel)

func NewDbArticleModel(author *models.Author, title string, ops ...ArticleDatabaseModelOpt) *ArticleDatabaseModel {
	m := &ArticleDatabaseModel{
		Title:      title,
		Content:    "",
		AuthorName: author.Name,
	}
	for _, opt := range ops {
		opt(m)
	}
	return m
}

func (a *ArticleDatabaseModel) ToDomain() *models.Article {
	return models.NewArticleOpt(
		models.ArticleWithTitle(a.Title),
		models.ArticleWithAuthor(a.AuthorName),
		models.ArticleWithId(a.ID),
		models.ArticleWithContent(a.Content))
}

func DbArticleModelWithContent(content string) ArticleDatabaseModelOpt {
	return func(m *ArticleDatabaseModel) {
		m.Content = content
	}
}

type DbArticleRepository struct {
	db *gorm.DB
}

var _ ports.ArticleRepository = (*DbArticleRepository)(nil)

func NewDbArticleRepository(db *gorm.DB) *DbArticleRepository {
	return &DbArticleRepository{db: db}
}

func (r *DbArticleRepository) Get(id uint) (*models.Article, error) {
	var db_article ArticleDatabaseModel
	if err := r.db.First(&db_article, id).Error; err != nil {
		return nil, err
	}
	return db_article.ToDomain(), nil
}

func (r *DbArticleRepository) Save(author *models.Author, title string, content string) (*models.Article, error) {
	db_model := NewDbArticleModel(author, title, DbArticleModelWithContent(content))
	if err := r.db.Create(&db_model).Error; err != nil {
		return nil, err
	}
	return db_model.ToDomain(), nil
}
