package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v5"
	"github.com/stretchr/testify/assert"
	"github.com/yaredow/new-arch/internal/repository/postgres"
)

func TestAuthorGetByID(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mockPool.Close()

	rows := mockPool.NewRows([]string{"id", "name", "updated_at", "created_at"})
	rows.AddRow(int64(1), "Iman", time.Now(), time.Now())
	mockPool.ExpectQuery("SELECT id, name, updated_at, created_at FROM author WHERE id = \\$1").WithArgs(int64(1)).WillReturnRows(rows)

	repo := postgres.NewAuthorRepository(mockPool)
	author, err := repo.GetByID(context.TODO(), 1)

	assert.NoError(t, err)
	assert.Equal(t, "Iman", author.Name)
	assert.NoError(t, mockPool.ExpectationsWereMet())
}
