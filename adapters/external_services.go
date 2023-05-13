package adapters

import (
	"fmt"

	"github.com/chernyshev-alex/go-hexagon/domain/models"
	"github.com/chernyshev-alex/go-hexagon/domain/ports"
)

type AuthorExternalModel struct {
	id                  string
	firstName, lastName string
}

func NewAuthorExternalModel(id, fName, lName string) *AuthorExternalModel {
	return &AuthorExternalModel{
		id:        id,
		firstName: fName,
		lastName:  lName,
	}
}

func (m AuthorExternalModel) FullName() string {
	return fmt.Sprintf("%s %s", m.firstName, m.lastName)
}

func (m AuthorExternalModel) ToDomain() *models.Author {
	return models.NewAuthor(m.id, m.FullName())

}

func (m AuthorExternalModel) ToString() string {
	return m.FullName()
}

type ExtSvcClientAuthorRepository struct {
	client interface{}
}

var _ ports.AuthorRepository = (*ExtSvcClientAuthorRepository)(nil)

func NewExtSvcClientAuthorRepository(client interface{}) *ExtSvcClientAuthorRepository {
	return &ExtSvcClientAuthorRepository{
		client: client,
	}
}

func (r ExtSvcClientAuthorRepository) Get(authorId string) (*models.Author, error) {
	return models.NewAuthor("1", "John Doyle"), nil
}
