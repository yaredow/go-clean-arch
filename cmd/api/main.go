package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/yaredow/new-arch/internal/handler"
	"github.com/yaredow/new-arch/internal/repository/postgres"
	"github.com/yaredow/new-arch/internal/service/article"
)

func main() {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASS"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"),
	)
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal("unable to connect to the database", err)
	}
	defer pool.Close()

	// Repos
	articleRepo := postgres.NewArticleRepository(pool)
	authRepo := postgres.NewAuthorRepository(pool)

	// Services
	svc := article.NewService(articleRepo, authRepo)

	// Echo
	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	// Handlers
	handler.NewArticleHandler(e, svc)

	e.GET("/", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Hello, World!"})
	})

	if err := e.Start(":4000"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
