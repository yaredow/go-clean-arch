package article

import (
	"context"
	"time"

	"github.com/yaredow/new-arch/internal/domain"
)

type ArticleRepository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []domain.Article, nextCursor string, err error)
	GetByID(ctx context.Context, id int64) (domain.Article, error)
	GetByTitle(ctx context.Context, title string) (domain.Article, error)
	Update(ctx context.Context, ar *domain.Article) error
	Store(ctx context.Context, a *domain.Article) error
	Delete(ctx context.Context, id int64) error
}

type AuthorRepository interface {
	GetByID(ctx context.Context, id int64) (domain.Author, error)
}

type Service struct {
	articleRepo ArticleRepository
	authorRepo  AuthorRepository
}

func NewService(a ArticleRepository, ar AuthorRepository) *Service {
	return &Service{
		articleRepo: a,
		authorRepo:  ar,
	}
}

func (a *Service) Fetch(ctx context.Context, cursor string, num int64) (res []domain.Article, nextCursor string, err error) {
	return a.articleRepo.Fetch(ctx, cursor, num)
}

func (a *Service) GetByID(ctx context.Context, id int64) (res domain.Article, err error) {
	res, err = a.articleRepo.GetByID(ctx, id)
	if err != nil {
		return
	}

	resAuthor, err := a.authorRepo.GetByID(ctx, res.Author.ID)
	if err != nil {
		return domain.Article{}, err
	}

	res.Author = resAuthor
	return
}

func (a *Service) GetByTitle(ctx context.Context, title string) (res domain.Article, err error) {
	res, err = a.articleRepo.GetByTitle(ctx, title)
	if err != nil {
		return
	}

	resAuthor, err := a.authorRepo.GetByID(ctx, res.Author.ID)
	if err != nil {
		return domain.Article{}, err
	}

	res.Author = resAuthor
	return
}

func (a *Service) Update(ctx context.Context, ar *domain.Article) (err error) {
	ar.UpdatedAt = time.Now()
	return a.articleRepo.Update(ctx, ar)
}

func (a *Service) Store(ctx context.Context, ar *domain.Article) (err error) {
	ar.CreatedAt = time.Now()
	ar.UpdatedAt = time.Now()
	return a.articleRepo.Store(ctx, ar)
}

func (a *Service) Delete(ctx context.Context, id int64) (err error) {
	existedArticle, err := a.articleRepo.GetByID(ctx, id)
	if err != nil {
		return
	}

	if existedArticle == (domain.Article{}) {
		return domain.ErrNotFound
	}

	return a.articleRepo.Delete(ctx, id)
}
