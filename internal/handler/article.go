package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
	"github.com/yaredow/new-arch/internal/domain"
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
	e.POST("/articles", h.Store)
	e.GET("/articles/:id", h.GetByID)
	e.PUT("/articles/:id", h.Update)
	e.DELETE("/articles/:id/delete", h.Delete)
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

func (h *ArticleHandler) GetByID(c *echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil || idP == 0 {
		return c.JSON(getStatusCode(err), map[string]string{"message": "not found"})
	}

	ctx := c.Request().Context()
	art, err := h.svc.GetByID(ctx, int64(idP))
	if err != nil {
		return c.JSON(getStatusCode(err), map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, art)
}

func (h *ArticleHandler) Store(c *echo.Context) error {
	var article domain.Article
	if err := c.Bind(&article); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"message": err.Error()})
	}

	if err := h.svc.Store(c.Request().Context(), &article); err != nil {
		return c.JSON(getStatusCode(err), map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, article)
}

func (h *ArticleHandler) Delete(c *echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "not found"})
	}

	if err := h.svc.Delete(c.Request().Context(), int64(idP)); err != nil {
		return c.JSON(getStatusCode(err), map[string]string{"message": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *ArticleHandler) Update(c *echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "not found"})
	}

	var article domain.Article
	if err = c.Bind(&article); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"message": err.Error()})
	}

	article.ID = int64(idP)
	if err := h.svc.Update(c.Request().Context(), &article); err != nil {
		return c.JSON(getStatusCode(err), map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, article)
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
