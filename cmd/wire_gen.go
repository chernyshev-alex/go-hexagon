// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package cmd

import (
	"github.com/chernyshev-alex/go-hexagon/internal/adapters"
	"github.com/chernyshev-alex/go-hexagon/internal/adapters/flows"
	"github.com/chernyshev-alex/go-hexagon/internal/domain/ports"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
)

import (
	_ "github.com/google/wire/cmd/wire"
)

// Injectors from wire.go:

func initializeWire() (*fiber.App, error) {
	db, err := initializeDB()
	if err != nil {
		return nil, err
	}
	articleRepository := initializeArticleRepository(db)
	authorRepository := initializeAuthorRepository()
	articleMessageSender := initializeArticleMessageSender()
	publishProviderWorkflow := initializePublishProviderWorkflow()
	articlePublisher := initializeArticlePublisher(articleMessageSender, publishProviderWorkflow)
	articleService := initializeArticleService(articleRepository, authorRepository, articlePublisher)
	articleFacade := initializeArticleFacade(articleService)
	fiberEndpoint := initializeEndpoint(articleFacade)
	app := initializeApp(fiberEndpoint)
	return app, nil
}

// wire.go:

func initializeDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db, nil
}

func initializeArticleRepository(db *gorm.DB) ports.ArticleRepository {
	return adapters.NewDbArticleRepository(db)
}

func initializeAuthorRepository() ports.AuthorRepository {
	return adapters.NewExtSvcClientAuthorRepository(http.DefaultClient)
}

func initializeArticleMessageSender() ports.ArticleMessageSender {
	return adapters.NewMbArticleMessageSender(http.DefaultClient)
}

func initializePublishProviderWorkflow() flows.PublishProviderWorkflow {
	return flows.PublishProviderWorkflow{}
}

func initializeArticlePublisher(articleMessageSender ports.ArticleMessageSender,
	publishProvider flows.PublishProviderWorkflow) *ports.ArticlePublisher {
	return ports.NewArticlePublisher(articleMessageSender, publishProvider)
}

func initializeArticleService(articleRepository ports.ArticleRepository,
	authorRepository ports.AuthorRepository,
	articlePublisher *ports.ArticlePublisher) *ports.ArticleService {
	return ports.NewArticleService(articleRepository, authorRepository, articlePublisher)
}

func initializeArticleFacade(articleService *ports.ArticleService) adapters.ArticleFacade {
	return adapters.NewArticleFacade(articleService)
}

func initializeEndpoint(facade adapters.ArticleFacade) adapters.HttpEndpoint {
	return adapters.NewEndpoint(facade)
}

func initializeApp(endpoint adapters.HttpEndpoint) *fiber.App {
	app := fiber.New()
	endpoint.AddRoutes(app)
	return app
}
