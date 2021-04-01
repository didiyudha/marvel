package handler

import (
	"context"
	"encoding/json"
	"github.com/didiyudha/marvel/business/data/character"
	"github.com/didiyudha/marvel/business/usecase/mock"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewHTTPHandler(t *testing.T) {
	mockctl := gomock.NewController(t)
	defer mockctl.Finish()

	marvelUseCase := mock.NewMockMarvelUseCase(mockctl)
	h := NewHTTPHandler(marvelUseCase)

	assert.NotNil(t, h)
	assert.NotNil(t, h.MarvelUseCase)
}

func TestGetAllCharacterID(t *testing.T) {
	mockctl := gomock.NewController(t)
	defer mockctl.Finish()

	t.Run("Successfully got all character ids", func(t *testing.T) {

		ctx := context.Background()

		req := httptest.NewRequest(http.MethodGet, "/characters", http.NoBody)
		req.WithContext(ctx)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)

		ids := []int{1, 2, 3, 4, 5}


		marvelUseCase := mock.NewMockMarvelUseCase(mockctl)
		marvelUseCase.
			EXPECT().
			GetAllCharacterID(ctx).
			Return(ids, nil)
		h := NewHTTPHandler(marvelUseCase)
		h.GetAllCharacterID(c)

		var recordID []int
		err := json.NewDecoder(rec.Body).Decode(&recordID)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, ids, recordID)

	})

	t.Run("Unsuccessful get all character id", func(t *testing.T) {

		ctx := context.Background()

		req := httptest.NewRequest(http.MethodGet, "/characters", http.NoBody)
		req.WithContext(ctx)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)



		marvelUseCase := mock.NewMockMarvelUseCase(mockctl)

		expectedErr := errors.New("database down")

		marvelUseCase.
			EXPECT().
			GetAllCharacterID(ctx).
			Return(nil, expectedErr)

		h := NewHTTPHandler(marvelUseCase)
		h.GetAllCharacterID(c)

		expectedErrResp := struct {
			Message string `json:"message"`
		}{}
		err := json.NewDecoder(rec.Body).Decode(&expectedErrResp)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, expectedErr.Error(), expectedErrResp.Message)

	})

}

func TestFindOne(t *testing.T) {
	mockctl := gomock.NewController(t)
	defer mockctl.Finish()

	t.Run("Successfully got a character of marvel", func(t *testing.T) {
		ctx := context.Background()

		req := httptest.NewRequest(http.MethodGet, "/characters/1", http.NoBody)
		req.WithContext(ctx)

		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		marvelUseCase := mock.NewMockMarvelUseCase(mockctl)
		characterMarvel := character.Character{
			ID:          1,
			Name:        "Ant Man",
			Description: "Awesome",
		}
		marvelUseCase.
			EXPECT().
			GetCharacter(ctx, 1).
			Return(characterMarvel, nil)
		h := NewHTTPHandler(marvelUseCase)
		h.FindOne(c)

		var expectedChar character.Character
		err := json.NewDecoder(rec.Body).Decode(&expectedChar)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expectedChar.ID, characterMarvel.ID)
		assert.Equal(t, expectedChar.Name, characterMarvel.Name)
		assert.Equal(t, expectedChar.Description, characterMarvel.Description)
	})

	t.Run("Character not found", func(t *testing.T) {
		ctx := context.Background()

		req := httptest.NewRequest(http.MethodGet, "/characters/1", http.NoBody)
		req.WithContext(ctx)

		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		marvelUseCase := mock.NewMockMarvelUseCase(mockctl)
		expectedErr := errors.Wrap(character.ErrNotFound, "get one character")
		marvelUseCase.
			EXPECT().
			GetCharacter(ctx, 1).
			Return(character.Character{}, expectedErr)

		h := NewHTTPHandler(marvelUseCase)
		h.FindOne(c)

		expectedErrResp := struct {
			Message string `json:"message"`
		}{}

		err := json.NewDecoder(rec.Body).Decode(&expectedErrResp)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, "character not found", expectedErrResp.Message)
	})

	t.Run("Unsuccessful to get a character of marvel", func(t *testing.T) {
		ctx := context.Background()

		req := httptest.NewRequest(http.MethodGet, "/characters/1", http.NoBody)
		req.WithContext(ctx)

		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		marvelUseCase := mock.NewMockMarvelUseCase(mockctl)
		expectedErr := errors.New("database down")
		marvelUseCase.
			EXPECT().
			GetCharacter(ctx, 1).
			Return(character.Character{}, expectedErr)
		h := NewHTTPHandler(marvelUseCase)
		h.FindOne(c)

		expectedErrResp := struct {
			Message string `json:"message"`
		}{}

		err := json.NewDecoder(rec.Body).Decode(&expectedErrResp)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, expectedErr.Error(), expectedErrResp.Message)

	})
}