package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yaredow/new-arch/domain"
)

type ArticleRepository struct {
	pool *pgxpool.Pool
}

func NewArticleRepository(pool *pgxpool.Pool) *ArticleRepository {
	return &ArticleRepository{pool: pool}
}

func (r *ArticleRepository) Fetch(ctx context.Context, cursor string, num int64) (res []domain.Article, nextCursor string, err error) {
	query := `
		SELECT id, title, content, author_id, updated_at, created_at
		FROM article WHERE created_at > $1 ORDER BY created_at LIMIT $2
	`

	var decodedCursor time.Time
	if cursor != "" {
		decodedCursor, err = decodeCursor(cursor)
		if err != nil {
			return nil, "", domain.ErrBadParamInput
		}
	}

	rows, err := r.pool.Query(ctx, query, decodedCursor, num)
	if err != nil {
		return nil, "", err
	}

	defer rows.Close()

	for rows.Next() {
		var a domain.Article
		var authorID int64
		err := rows.Scan(&a.ID, &a.Title, &a.Content, &authorID, &a.UpdatedAt, &a.CreatedAt)
		if err != nil {
			return nil, "", err
		}

		a.Author = domain.Author{ID: authorID}
		res = append(res, a)
	}

	if len(res) == int(num) {
		nextCursor = encodeCursor(res[len(res)-1].CreatedAt)
	}
	return
}

func (r *ArticleRepository) GetByID(ctx context.Context, id int64) (domain.Article, error) {
	query := `
		SELECT id, title, content, author_id, updated_at, created_at FROM articles WHERE id = $1
	`
	var a domain.Article
	var authorID int64
	err := r.pool.QueryRow(ctx, query, id).Scan(&a.ID, &a.Title, &a.Content, &authorID, &a.UpdatedAt, &a.CreatedAt)
	if err != nil {
		switch {
		case err.Error() == "no rows in result set":
			return domain.Article{}, domain.ErrNotFound
		default:
			return domain.Article{}, err
		}
	}
	a.Author.ID = authorID
	return a, nil
}

func (r *ArticleRepository) GetByTitle(ctx context.Context, title string) (domain.Article, error) {
	query := `
		SELECT id, title, content, author_id, updated_at, created_at FROM articles WHERE title = $1
	`

	var a domain.Article
	var authorID int64
	err := r.pool.QueryRow(ctx, query, title).Scan(
		&a.ID, &a.Title, &a.Content, &authorID, &a.UpdatedAt, &a.CreatedAt,
	)
	if err != nil {
		switch {
		case err.Error() == "no rows in result set":
			return domain.Article{}, domain.ErrNotFound
		default:
			return domain.Article{}, err
		}
	}

	a.Author.ID = authorID
	return a, nil
}

func (r *ArticleRepository) Store(ctx context.Context, a *domain.Article) error {
	query := `
		INSERT INTO articles (title, content, author_id, updated_at, created_at) VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	args := []any{a.Title, a.Content, a.Author.ID, a.UpdatedAt, a.CreatedAt}
	err := r.pool.QueryRow(ctx, query, args...).Scan(&a.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *ArticleRepository) Update(ctx context.Context, a *domain.Article) (err error) {
	query := `
		UPDATE article SET title = $1, content = $2, author_id = $3, updated_at = $4 WHERE id = $5
	`

	args := []any{a.Title, a.Content, a.Author.ID, a.UpdatedAt, a.ID}
	tag, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *ArticleRepository) Delete(ctx context.Context, id int64) (err error) {
	query := `
		DELETE FROM articles WHERE id = $1
	`

	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}
