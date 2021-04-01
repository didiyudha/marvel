package api

import (
	"fmt"
	"github.com/didiyudha/marvel/business/data/character"
	"github.com/didiyudha/marvel/business/usecase"
	"github.com/didiyudha/marvel/business/web/handler"
	"github.com/didiyudha/marvel/business/web/routes"
	"github.com/didiyudha/marvel/foundation/database"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"net/http"
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
	return a.e.Start(fmt.Sprintf(":%d", port))
}
