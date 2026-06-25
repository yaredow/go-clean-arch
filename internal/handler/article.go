package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
	"github.com/yaredow/new-arch/domain"
)

type ArticleService interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []domain.Article, nextCursor string, err error)
	GetByID(ctx context.Context, id int64) (domain.Article, error)
	GetByTitle(ctx context.Context, title string) (domain.Article, error)
	Update(ctx context.Context, ar *domain.Article) error
	Store(ctx context.Context, a *domain.Article) error
	Delete(ctx context.Context, id int64) error
}

type ArticleHandler struct {
	svc ArticleService
}

func NewArticleHandler(e *echo.Echo, svc ArticleService) {
	h := &ArticleHandler{svc: svc}
	e.GET("/articles", h.FetchArticle)
}

func (h *ArticleHandler) FetchArticle(c *echo.Context) error {
	numS := c.QueryParam("num")
	num, err := strconv.Atoi(numS)
	if err != nil || num == 0 {
		num = 10
	}

	cursor := c.QueryParam("cursor")
	ctx := c.Request().Context()

	listAr, nextCursor, err := h.svc.Fetch(ctx, cursor, int64(num))
	if err != nil {
		return c.JSON(getStatusCode(err), map[string]string{"message": err.Error()})
	}

	c.Response().Header().Set("X-Cursor", nextCursor)
	return c.JSON(http.StatusOK, listAr)
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	case domain.ErrBadParamInput:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
