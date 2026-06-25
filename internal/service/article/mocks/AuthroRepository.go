package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/yaredow/new-arch/internal/domain"
)

type AuthorRepository struct {
	mock.Mock
}

func (m *AuthorRepository) GetByID(ctx context.Context, id int64) (domain.Author, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Author), args.Error(1)
}
