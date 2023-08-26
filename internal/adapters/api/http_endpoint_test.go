package api_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chernyshev-alex/go-hexagon/internal/adapters/api"
	"github.com/chernyshev-alex/go-hexagon/internal/adapters/api/mocks"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetArticle(t *testing.T) {
	mf := mocks.NewArticleFacade(t)

	mf.On("Get", mock.Anything, "1111").Return(api.ArticleResponse{ID: 1111, Title: "title", AuthorName: "some author"}, nil)
	mf.On("Get", mock.Anything, "0000").Return(api.ArticleResponse{}, fmt.Errorf("not found"))

	app := configureEndPoint(mf)

	req := httptest.NewRequest(http.MethodGet, "/articles/1111", nil)
	res, _ := app.Test(req, -1)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	article, err := decodeArticleResponse(res.Body)

	assert.Nil(t, err)
	assert.Equal(t, "1111", article.ID)

	req = httptest.NewRequest(http.MethodGet, "/articles/0000", nil)
	if res, err = app.Test(req, -1); err != nil {
		fmt.Println(err)
	}
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

/*
func TestSearchArticle(t *testing.T) {
	facade := MockedFacade{}
	facade.On("SearchBy", "authorname", "Conan Doyle").Return(
		[]*ArticleResponse{
			{
				Id:         "1",
				Title:      "ATitle",
				Content:    "",
				AuthorName: "Conan Doyle",
			},
		}, nil)

	var q = `query { articles(title:"ATitle", authorname:"Conan Doyle") {
				id title content authorname
			}}`

	app := configureEndPoint(&facade)
	req := newJsonRequest(http.MethodPost, "/search/articles", gqlRequestBody{Query: q})

	res, _ := app.Test(req, -1)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var gqlResponse struct {
		Data struct {
			Articles []*ArticleResponse `json:"articles"`
		} `json:"data"`
	}
	if err := json.NewDecoder(res.Body).Decode(&gqlResponse); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(gqlResponse.Data.Articles), 1)
	assert.Equal(t, gqlResponse.Data.Articles[0].AuthorName, "Conan Doyle")
}

func TestCreateArticle(t *testing.T) {
	facade := MockedFacade{}
	request := ArticleRequest{Title: "title", Content: "content", AuthorId: "authorId"}

	facade.On("Create", &request).Return(ArticleIdResponse{Id: "1"}, nil)

	app := configureEndPoint(&facade)
	req := newJsonRequest(http.MethodPost, "/articles", request)

	res, _ := app.Test(req, -1)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	article, err := decodeArticleResponse(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, "1", article.Id)

} */

func configureEndPoint(af *mocks.ArticleFacade) *fiber.App {
	endpoint := api.NewEndpoint(af)
	app := fiber.New()
	endpoint.AddRoutes(app)
	return app
}

// func newJsonRequest(httpMethod, Url string, v interface{}) *http.Request {
// 	body, _ := json.Marshal(v)
// 	req := httptest.NewRequest(httpMethod, Url, bytes.NewReader(body))
// 	req.Header.Set("Content-Type", "application/json")
// 	return req
// }

func decodeArticleResponse(r io.ReadCloser) (api.ArticleResponse, error) {
	var article api.ArticleResponse
	err := json.NewDecoder(r).Decode(&article)
	if err != nil {
		return api.ArticleResponse{}, fmt.Errorf("error %v", err)
	}
	return article, nil
}
