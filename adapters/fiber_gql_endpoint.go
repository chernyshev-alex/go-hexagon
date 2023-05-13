package adapters

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/graphql-go/graphql"
)

const (
	gql_AuthorName = "authorname"
	gql_Title      = "title"
	gql_Id         = "id"
	gql_Content    = "content"
)

type gqlRequestBody struct {
	Query string `json:"query"`
}

func (ap *FiberEndpoint) SearchBy(what, value string) ([]*ArticleResponse, error) {
	return ap.articleFacade.SearchBy(what, value)
}

var articleType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "article",
		Fields: graphql.Fields{
			gql_Id:         &graphql.Field{Type: graphql.String},
			gql_Title:      &graphql.Field{Type: graphql.String},
			gql_Content:    &graphql.Field{Type: graphql.String},
			gql_AuthorName: &graphql.Field{Type: graphql.String},
		},
	})

func (ap *FiberEndpoint) FieldsResolver(params graphql.ResolveParams) (interface{}, error) {
	var articles []*ArticleResponse
	title, _ := params.Args[gql_Title].(string)
	authorName, _ := params.Args[gql_AuthorName].(string)

	if len(authorName) > 0 {
		result, err := ap.articleFacade.SearchBy(gql_AuthorName, authorName)
		if err != nil {
			return nil, err
		}
		for _, r := range result {
			if len(title) > 0 {
				if title == r.Title {
					articles = append(articles, r)
				}
			} else {
				articles = append(articles, r)
			}
		}
	}
	return articles, nil
}

func (ap *FiberEndpoint) getGraphQlHandler() fiber.Handler {
	var queryType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",

		Fields: graphql.Fields{
			"articles": &graphql.Field{
				Type: graphql.NewList(articleType),
				Args: graphql.FieldConfigArgument{
					gql_Title:      &graphql.ArgumentConfig{Type: graphql.String},
					gql_AuthorName: &graphql.ArgumentConfig{Type: graphql.String},
				},

				Resolve: ap.FieldsResolver,
			},
		}})

	return func(c *fiber.Ctx) error {
		var reqBody gqlRequestBody
		if err := json.Unmarshal([]byte(c.Body()), &reqBody); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		schema, _ := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: reqBody.Query,
		})

		if result.HasErrors() {
			return c.Status(fiber.StatusBadRequest).JSON(result.Errors)
		}
		return c.JSON(result)
	}
}
