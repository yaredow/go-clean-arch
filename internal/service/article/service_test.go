package article_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yaredow/new-arch/internal/domain"
	"github.com/yaredow/new-arch/internal/service/article"
	"github.com/yaredow/new-arch/internal/service/article/mocks"
)

func TestFetchArticle(t *testing.T) {
	mockArticle := domain.Article{
		Title:   "Hello",
		Content: "Content",
	}
	mockList := []domain.Article{mockArticle}

	t.Run("success", func(t *testing.T) {
		mockArticleRepo := new(mocks.ArticleRepository)
		mockAuthorRepo := new(mocks.AuthorRepository)

		mockArticleRepo.On("Fetch", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int64")).Return(mockList, "next-cursor", nil).Once()

		svc := article.NewService(mockArticleRepo, mockAuthorRepo)
		list, nextCursor, err := svc.Fetch(context.TODO(), "12", 1)

		assert.Equal(t, "next-cursor", nextCursor)
		assert.NoError(t, err)
		assert.Len(t, list, 1)
		mockArticleRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockArticleRepo := new(mocks.ArticleRepository)
		mockAuthorRepo := new(mocks.AuthorRepository)

		mockArticleRepo.On("Fetch", mock.Anything, mock.AnythingOfType("string"),
			mock.AnythingOfType("int64")).Return(nil, "", assert.AnError).Once()

		svc := article.NewService(mockArticleRepo, mockAuthorRepo)
		list, nextCursor, err := svc.Fetch(context.TODO(), "12", 1)

		assert.Empty(t, nextCursor)
		assert.Error(t, err)
		assert.Len(t, list, 0)
		mockArticleRepo.AssertExpectations(t)
	})
}
