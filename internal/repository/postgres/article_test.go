package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v5"
	"github.com/stretchr/testify/assert"
	"github.com/yaredow/new-arch/internal/repository/postgres"
)

func TestGetByID(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mockPool.Close()

	now := time.Now()
	rows := mockPool.NewRows([]string{"id", "title", "content", "author_id", "updated_at", "created_at"}).
		AddRow(1, "Title", "Content", 1, now, now)

	mockPool.ExpectQuery("SELECT id, title, content, author_id, updated_at, created_at FROM article WHERE id = \\$1").
		WithArgs(int64(1)).WillReturnRows(rows)

	repo := postgres.NewArticleRepository(mockPool)
	article, err := repo.GetByID(context.TODO(), 1)

	assert.NoError(t, err)
	assert.Equal(t, "Title", article.Title)
	assert.Equal(t, int64(1), article.Author.ID)

	assert.NoError(t, mockPool.ExpectationsWereMet())
}
