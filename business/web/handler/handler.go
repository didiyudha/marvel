package handler

import (
	"github.com/didiyudha/marvel/business/data/character"
	"github.com/didiyudha/marvel/business/usecase"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

type HTTPHandler struct {
	MarvelUseCase usecase.MarvelUseCase
}

func NewHTTPHandler(marvelUseCase usecase.MarvelUseCase) *HTTPHandler {
	return &HTTPHandler{MarvelUseCase: marvelUseCase}
}

func (h *HTTPHandler) GetAllCharacterID(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := h.MarvelUseCase.GetAllCharacterID(ctx)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, id)
}

func (h *HTTPHandler) FindOne(c echo.Context) error {

	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid character id")
	}

	char, err := h.MarvelUseCase.GetCharacter(ctx, id)
	if errors.Cause(err) == character.ErrNotFound {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"message": "character not found",
		})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, char)
}