package routes

import (
	"github.com/didiyudha/marvel/business/web/handler"
	"github.com/labstack/echo/v4"
)

func Mount(e *echo.Echo, h *handler.HTTPHandler) {
	e.GET("/characters", h.GetAllCharacterID)
	e.GET("/characters/:id", h.FindOne)
}
