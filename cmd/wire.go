//go:build wireinject
// +build wireinject

package cmd

import (
	"net/http"

	"github.com/chernyshev-alex/go-hexagon/internal/adapters"
	"github.com/chernyshev-alex/go-hexagon/internal/adapters/flows"
	"github.com/chernyshev-alex/go-hexagon/internal/domain/ports"
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	_ "github.com/google/wire/cmd/wire"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

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

func initializeEndpoint(facade adapters.ArticleFacade) adapters.FiberEndpoint {
	return adapters.NewEndpoint(facade)
}

func initializeApp(endpoint adapters.FiberEndpoint) *fiber.App {
	app := fiber.New()
	endpoint.AddRoutes(app)
	return app
}

func initializeWire() (*fiber.App, error) {
	wire.Build(
		initializeDB,
		initializeArticleRepository,
		initializeAuthorRepository,
		initializeArticleMessageSender,
		initializePublishProviderWorkflow,
		initializeArticlePublisher,
		initializeArticleService,
		initializeArticleFacade,
		initializeEndpoint,
		initializeApp,
	)
	return nil, nil
}
