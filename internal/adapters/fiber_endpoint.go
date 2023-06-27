package adapters

import (
	"github.com/gofiber/fiber/v2"
)

type FiberEndpoint struct {
	articleFacade ArticleFacade
}

func NewEndpoint(af ArticleFacade) FiberEndpoint {
	return FiberEndpoint{articleFacade: af}
}

func (ap *FiberEndpoint) AddRoutes(app *fiber.App) {
	app.Get("/articles/:id", ap.GetArticle)
	app.Post("/articles", ap.CreateArticle)
	app.Post("/search/articles", ap.getGraphQlHandler())
}

func (ap *FiberEndpoint) GetArticle(c *fiber.Ctx) error {
	if articleResponse, err := ap.articleFacade.Get(c.Params("id")); err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	} else {
		return c.JSON(articleResponse)
	}
}

func (ap *FiberEndpoint) CreateArticle(c *fiber.Ctx) error {
	var req ArticleRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if res, err := ap.articleFacade.Create(&req); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else {
		return c.JSON(res)
	}
}
