package api

import (
	"github.com/gofiber/fiber/v2"
)

type HttpEndpoint struct {
	articleFacade ArticleFacade
}

func NewEndpoint(af ArticleFacade) HttpEndpoint {
	return HttpEndpoint{articleFacade: af}
}

func (ap *HttpEndpoint) AddRoutes(app *fiber.App) {
	app.Get("/articles/:id", ap.GetArticle)
	app.Post("/articles", ap.CreateArticle)
	app.Post("/search/articles", ap.getGraphQlHandler())
}

func (ap *HttpEndpoint) GetArticle(c *fiber.Ctx) error {
	if articleResponse, err := ap.articleFacade.Get(c.Context(), c.Params("id")); err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	} else {
		return c.JSON(articleResponse)
	}
}

func (ap *HttpEndpoint) CreateArticle(c *fiber.Ctx) error {
	var req ArticleRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if res, err := ap.articleFacade.Create(c.Context(), req); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else {
		return c.JSON(res)
	}
}
