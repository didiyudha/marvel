package api

import (
	"context"
	"fmt"
	"github.com/didiyudha/marvel/business/data/character"
	"github.com/didiyudha/marvel/business/usecase"
	"github.com/didiyudha/marvel/business/web/handler"
	"github.com/didiyudha/marvel/business/web/routes"
	"github.com/didiyudha/marvel/foundation/database"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Api struct {
	e *echo.Echo
}

func New(db *sqlx.DB, redisClient *redis.Client) *Api {

	characterStore := character.NewStore(db)
	characterCaching := character.NewCaching(redisClient)
	marvelUseCase := usecase.NewMarvelUseCase(characterStore, characterCaching)
	httpHandler := handler.NewHTTPHandler(marvelUseCase)

	e := echo.New()
	e.Use(middleware.Logger())
	e.GET("/healthy", healthy(db))
	routes.Mount(e, httpHandler)

	return &Api{e: e}
}

func healthy(db *sqlx.DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		if err := database.StatusCheck(ctx, db); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": err.Error(),
			})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": http.StatusText(http.StatusOK),
		})
	}
}

func (a *Api) Serve(port int) error {

	go func() {
		if err := a.e.Start(fmt.Sprintf(":%d", port)); err != nil && err != http.ErrServerClosed {
			log.Printf("Error : %v\n", err)
			log.Printf("Server shutting down...")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := a.e.Shutdown(ctx); err != nil {
		a.e.Logger.Fatal(err)
	}

	return nil
}
