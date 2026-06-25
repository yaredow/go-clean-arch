package postgres

import (
	"context"

	"github.com/yaredow/new-arch/internal/domain"
)

type AuthorRepository struct {
	pool Pool
}

func NewAuthorRepository(pool Pool) *AuthorRepository {
	return &AuthorRepository{pool: pool}
}

func (r *AuthorRepository) GetByID(ctx context.Context, id int64) (domain.Author, error) {
	query := `SELECT id, name, updated_at, created_at FROM author WHERE id = $1`

	var a domain.Author
	err := r.pool.QueryRow(ctx, query, id).Scan(&a.ID, &a.Name, &a.UpdatedAt, &a.CreatedAt)
	if err != nil {
		switch {
		case err.Error() == "now rows in result set":
			return domain.Author{}, nil
		default:
			return domain.Author{}, err
		}
	}
	return a, nil
}
