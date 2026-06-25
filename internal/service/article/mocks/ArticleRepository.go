package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/yaredow/new-arch/internal/domain"
)

type ArticleRepository struct {
	mock.Mock
}

func (m *ArticleRepository) Fetch(ctx context.Context, cursor string, num int64) ([]domain.Article, string, error) {
	args := m.Called(ctx, cursor, num)
	
	res, ok := args.Get(0).([]domain.Article)
	if !ok {
		res = nil
	}
	
	return res, args.String(1), args.Error(2)
}

func (m *ArticleRepository) GetByID(ctx context.Context, id int64) (domain.Article, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Article), args.Error(1)
}

func (m *ArticleRepository) GetByTitle(ctx context.Context, title string) (domain.Article, error) {
	args := m.Called(ctx, title)
	return args.Get(0).(domain.Article), args.Error(1)
}

func (m *ArticleRepository) Update(ctx context.Context, ar *domain.Article) error {
	args := m.Called(ctx, ar)
	return args.Error(0)
}

func (m *ArticleRepository) Store(ctx context.Context, a *domain.Article) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m *ArticleRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
