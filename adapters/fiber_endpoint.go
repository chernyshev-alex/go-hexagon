package adapters

import (
	"github.com/gofiber/fiber/v2"
)

type FiberEndpoint struct {
	articleFacade ArticleFacade
}

func NewEndpoint(af ArticleFacade) FiberEndpoint {
	return FiberEndpoint{
		articleFacade: af,
	}
}

func (ap *FiberEndpoint) addRoutes(app *fiber.App) {
	app.Get("/articles/:id", ap.GetArticle)
	app.Post("/articles", ap.CreateArticle)
	app.Post("/search/articles", ap.getGraphQlHandler())
}

func (ap *FiberEndpoint) GetArticle(c *fiber.Ctx) error {
	articleId := c.Params("id")
	articleResponse, err := ap.articleFacade.Get(articleId)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	return c.JSON(articleResponse)
}

func (ap *FiberEndpoint) CreateArticle(c *fiber.Ctx) error {
	var req ArticleRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	res, err := ap.articleFacade.Create(&req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(res)
}
