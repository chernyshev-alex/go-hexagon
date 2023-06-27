package adapters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockedFacade struct {
	mock.Mock
}

var _ ArticleFacade = (*MockedFacade)(nil)

func (f MockedFacade) Get(articleId string) (*ArticleResponse, error) {
	args := f.Called(articleId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ArticleResponse), args.Error(1)
}

func (f MockedFacade) Create(rq *ArticleRequest) (ArticleIdResponse, error) {
	args := f.Called(rq)
	return args.Get(0).(ArticleIdResponse), args.Error(1)
}

func (f MockedFacade) List(pageId string) ([]*ArticleResponse, error) {
	args := f.Called(pageId)
	return args.Get(0).([]*ArticleResponse), args.Error(1)
}

func (f MockedFacade) SearchBy(what, value string) ([]*ArticleResponse, error) {
	args := f.Called(what, value)
	ww := args.Get(0)
	return ww.([]*ArticleResponse), args.Error(1)
}

func TestGetArticle(t *testing.T) {
	facade := MockedFacade{}
	facade.On("Get", "1").Return(&ArticleResponse{Id: "1", Title: "title-1"}, nil)
	facade.On("Get", "9999").Return(nil, fmt.Errorf("not found 999"))

	app := configureEndPoint(facade)

	req := httptest.NewRequest(http.MethodGet, "/articles/1", nil)
	res, _ := app.Test(req, -1)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	article, err := decodeArticleResponse(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, "1", article.Id)

	req = httptest.NewRequest(http.MethodGet, "/articles/9999", nil)
	res, _ = app.Test(req, -1)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

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

	app := configureEndPoint(facade)
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

	app := configureEndPoint(facade)
	req := newJsonRequest(http.MethodPost, "/articles", request)

	res, _ := app.Test(req, -1)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	article, err := decodeArticleResponse(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, "1", article.Id)

}

func configureEndPoint(facade MockedFacade) *fiber.App {
	endpoint := NewEndpoint(facade)
	app := fiber.New()
	endpoint.AddRoutes(app)
	return app
}

func newJsonRequest(httpMethod, Url string, v interface{}) *http.Request {
	body, _ := json.Marshal(v)
	req := httptest.NewRequest(httpMethod, Url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func decodeArticleResponse(r io.ReadCloser) (*ArticleResponse, error) {
	var article ArticleResponse
	err := json.NewDecoder(r).Decode(&article)
	if err != nil {
		return nil, fmt.Errorf("error %v", err)
	}
	return &article, nil
}
