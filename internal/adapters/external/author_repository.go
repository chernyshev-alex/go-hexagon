package adapters

import (
	"context"

	"github.com/chernyshev-alex/go-hexagon/internal/domain/model"
	"github.com/chernyshev-alex/go-hexagon/internal/domain/ports"
)

type AuthorExternalModel struct {
	id        int64
	firstName string
	lastName  string
}

func NewAuthorExternalModel(id int64, fname, lname string) AuthorExternalModel {
	return AuthorExternalModel{
		id:        id,
		firstName: fname,
		lastName:  lname,
	}
}

func (m AuthorExternalModel) ToDomain() model.Author {
	return model.Author{
		ID:   m.id,
		Name: m.lastName,
	}
}

type ExternalServiceClientAuthorRepository struct {
	ports.AuthorRepository
}

func NewExternalServiceClientAuthorRepository() ports.AuthorRepository {
	return ExternalServiceClientAuthorRepository{}
}

func (s ExternalServiceClientAuthorRepository) GetByID(ctx context.Context, id int64) (model.Author, error) {
	return NewAuthorExternalModel(id,
		"William",
		"Shakespeare").
		ToDomain(), nil
}
