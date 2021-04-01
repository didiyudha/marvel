package internal

import (
	"context"
	"github.com/bxcodec/faker/v3"
	"github.com/didiyudha/marvel/business/data/character"
	charactermock "github.com/didiyudha/marvel/business/data/character/mock"
	"github.com/didiyudha/marvel/client"
	clientmock "github.com/didiyudha/marvel/client/mock"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func populateCharacters(minID int, maxID int) []character.Character {
	characters := make([]character.Character, 0, maxID-minID)
	for i := minID; i <= maxID; i++ {
		c := character.Character{
			ID:          i,
			Name:        faker.Name(),
			Description: faker.Paragraph(),
		}
		characters = append(characters, c)
	}
	return characters
}

func populateMarvelCharacterResponse(total int) client.CharacterResponse {
	results := make([]client.Result, 0, total)
	for i := 1; i <= total; i++ {
		r := client.Result{
			ID:          i,
			Name:        faker.Name(),
			Description: faker.Paragraph(),
		}
		results = append(results, r)
	}
	resp := client.CharacterResponse{
		Data: client.Data{
			Results: results,
		},
	}
	return resp
}

func TestNewMarvelWorker(t *testing.T) {
	mockctl := gomock.NewController(t)
	defer mockctl.Finish()

	store := charactermock.NewMockStore(mockctl)
	marvelClient := clientmock.NewMockMarvelClient(mockctl)

	marvelWorker := NewMarvelWorker(store, marvelClient)
	assert.NotNil(t, marvelWorker)

	impl, ok := marvelWorker.(*marvelWorkerImpl)
	assert.True(t, ok)
	assert.NotNil(t, impl)

}

func TestFindNewData(t *testing.T) {
	t.Run("When there is new data", func(t *testing.T) {
		t.Run("It should return new data", func(t *testing.T) {
			fromDB := populateCharacters(1, 10)
			fromAPI := populateCharacters(1, 15)
			newCharacters := findNewData(fromAPI, fromDB)
			assert.Equal(t, 5, len(newCharacters))
			assert.Equal(t, 11, newCharacters[0].ID)
			assert.Equal(t, 12, newCharacters[1].ID)
			assert.Equal(t, 13, newCharacters[2].ID)
			assert.Equal(t, 14, newCharacters[3].ID)
			assert.Equal(t, 15, newCharacters[4].ID)
		})
	})
	t.Run("When there is no new data", func(t *testing.T) {
		t.Run("It should return empty new data", func(t *testing.T) {
			fromDB := populateCharacters(1, 10)
			fromAPI := populateCharacters(1, 10)
			newCharacters := findNewData(fromAPI, fromDB)
			assert.Equal(t, 0, len(newCharacters))
		})
	})
}

func TestGetTotalCharacter(t *testing.T) {
	mockctl := gomock.NewController(t)
	defer mockctl.Finish()

	ctx := context.Background()
	marvelClient := clientmock.NewMockMarvelClient(mockctl)
	t.Run("When error occurred while getting characters from marvel API", func(t *testing.T) {
		t.Run("It should return error with 0 total data", func(t *testing.T) {
			expectedErr := errors.New("")
			marvelClient.
				EXPECT().
				Characters(ctx, gomock.AssignableToTypeOf(time.Time{}), 1, 0).
				Return(client.CharacterResponse{}, expectedErr)
			m := &marvelWorkerImpl{
				MarvelClient: marvelClient,
			}
			total, err := m.GetTotalCharacter(ctx)
			assert.Equal(t, 0, total)
			assert.True(t, errors.Cause(err) == expectedErr)
		})
	})

	t.Run("When successfully got characters from marvel API", func(t *testing.T) {
		t.Run("It should return correct total data and no error", func(t *testing.T) {

			const totalData = 10

			resp := client.CharacterResponse{
				Data: client.Data{
					Total: totalData,
				},
			}

			marvelClient.
				EXPECT().
				Characters(ctx, gomock.AssignableToTypeOf(time.Time{}), 1, 0).
				Return(resp, nil)

			m := &marvelWorkerImpl{
				MarvelClient: marvelClient,
			}

			total, err := m.GetTotalCharacter(ctx)
			assert.NoError(t, err)
			assert.Equal(t, 10, total)
		})
	})
}

func TestGetAllCharactersFromAPI(t *testing.T) {
	mockctl := gomock.NewController(t)
	defer mockctl.Finish()

	ctx := context.Background()
	marvelClient := clientmock.NewMockMarvelClient(mockctl)

	t.Run("When error getting total character", func(t *testing.T) {
		t.Run("It should return empty data with error", func(t *testing.T) {

			expectedErr := errors.New("database down")

			m := &marvelWorkerImpl{
				MarvelClient: marvelClient,
			}

			marvelClient.
				EXPECT().
				Characters(ctx, gomock.AssignableToTypeOf(time.Time{}), 1, 0).
				Return(client.CharacterResponse{}, expectedErr)

			chars, err := m.GetAllCharactersFromAPI(ctx)
			assert.Nil(t, chars)
			assert.True(t, errors.Cause(err) == expectedErr)
		})
	})

	t.Run("When error occurred while getting characters from marvel api", func(t *testing.T) {
		t.Run("It should return empty data with error", func(t *testing.T) {
			expectedErr := errors.New("api down")
			const totalData = 10

			resp := client.CharacterResponse{
				Data: client.Data{
					Total: totalData,
				},
			}

			gomock.InOrder(
				marvelClient.
					EXPECT().
					Characters(ctx, gomock.AssignableToTypeOf(time.Time{}), 1, 0).
					Return(resp, nil),
				marvelClient.
					EXPECT().
					Characters(ctx, gomock.AssignableToTypeOf(time.Time{}), maxPerProcess, gomock.Any()).
					Return(client.CharacterResponse{}, expectedErr),
			)

			m := &marvelWorkerImpl{
				MarvelClient: marvelClient,
			}

			chars, err := m.GetAllCharactersFromAPI(ctx)
			assert.Nil(t, chars)
			assert.True(t, errors.Cause(err) == expectedErr)
		})
	})

	t.Run("When successfully got characters from api", func(t *testing.T) {
		t.Run("It should return characters data with no error", func(t *testing.T) {
			const totalData = 10
			resp := client.CharacterResponse{
				Data: client.Data{
					Total: totalData,
				},
			}
			gomock.InOrder(
				marvelClient.
					EXPECT().
					Characters(ctx, gomock.AssignableToTypeOf(time.Time{}), 1, 0).
					Return(resp, nil),
				marvelClient.
					EXPECT().
					Characters(ctx, gomock.AssignableToTypeOf(time.Time{}), maxPerProcess, gomock.Any()).
					Return(populateMarvelCharacterResponse(totalData), nil).
					Times(1),
			)
			m := &marvelWorkerImpl{
				MarvelClient: marvelClient,
			}
			chars, err := m.GetAllCharactersFromAPI(ctx)
			assert.NoError(t, err)
			assert.Equal(t, 10, len(chars))
			assert.Equal(t, 1, chars[0].ID)
			assert.Equal(t, 2, chars[1].ID)
			assert.Equal(t, 3, chars[2].ID)
			assert.Equal(t, 4, chars[3].ID)
			assert.Equal(t, 5, chars[4].ID)
			assert.Equal(t, 6, chars[5].ID)
			assert.Equal(t, 7, chars[6].ID)
			assert.Equal(t, 8, chars[7].ID)
			assert.Equal(t, 9, chars[8].ID)
			assert.Equal(t, 10, chars[9].ID)
		})
	})
}
