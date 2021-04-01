package exec

import (
	"context"
	"github.com/bxcodec/faker/v3"
	"github.com/didiyudha/marvel/business/data/character"
	charactermock "github.com/didiyudha/marvel/business/data/character/mock"
	"github.com/didiyudha/marvel/client"
	"github.com/didiyudha/marvel/client/mock"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

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

func TestInitializeMarvelCharacter(t *testing.T) {
	mockctl := gomock.NewController(t)
	defer mockctl.Finish()

	ctx := context.Background()
	store := charactermock.NewMockStore(mockctl)
	marvelClient := mock.NewMockMarvelClient(mockctl)

	m := &marvelCmdExecImpl{
		Store:        store,
		MarvelClient: marvelClient,
	}

	t.Run("When error occurred while clean up the data", func(t *testing.T) {
		t.Run("It should return error", func(t *testing.T) {
			expectedErr := errors.New("database down")
			store.
				EXPECT().
				CleanUp(ctx).
				Return(expectedErr)
			err := m.InitializeMarvelCharacter(ctx)
			assert.Error(t, err)
			assert.True(t, errors.Cause(err) == expectedErr)
		})
	})

	t.Run("When error occurred while get total data", func(t *testing.T) {
		t.Run("It should return error", func(t *testing.T) {
			expectedErr := errors.New("database down")
			gomock.InOrder(
				store.
					EXPECT().
					CleanUp(ctx).
					Return(nil),
				marvelClient.
					EXPECT().
					Characters(ctx, gomock.AssignableToTypeOf(time.Time{}), 1, 0).
					Return(client.CharacterResponse{}, expectedErr),
			)

			err := m.InitializeMarvelCharacter(ctx)
			assert.Error(t, err)
			assert.True(t, errors.Cause(err) == expectedErr)
		})
	})

	t.Run("When error occurred while getting characters from marvel api", func(t *testing.T) {
		t.Run("It should return error", func(t *testing.T) {
			expectedErr := errors.New("api down")
			const totalData = 10

			resp := client.CharacterResponse{
				Data: client.Data{
					Total: totalData,
				},
			}
			gomock.InOrder(
				store.
					EXPECT().
					CleanUp(ctx).
					Return(nil),
				marvelClient.
					EXPECT().
					Characters(ctx, gomock.AssignableToTypeOf(time.Time{}), 1, 0).
					Return(resp, nil),
				marvelClient.
					EXPECT().
					Characters(ctx, gomock.AssignableToTypeOf(time.Time{}), maxPerProcess, gomock.Any()).
					Return(client.CharacterResponse{}, expectedErr),
			)
			err := m.InitializeMarvelCharacter(ctx)
			assert.Error(t, err)
			assert.True(t, errors.Cause(err) == expectedErr)
		})
	})

	t.Run("When error occurred save characters data", func(t *testing.T) {
		t.Run("It should return error", func(t *testing.T) {

			expectedErr := errors.New("api down")
			const totalData = 10

			resp := client.CharacterResponse{
				Data: client.Data{
					Total: totalData,
				},
			}

			gomock.InOrder(
				store.
					EXPECT().
					CleanUp(ctx).
					Return(nil),
				marvelClient.
					EXPECT().
					Characters(ctx, gomock.AssignableToTypeOf(time.Time{}), 1, 0).
					Return(resp, nil),
				marvelClient.
					EXPECT().
					Characters(ctx, gomock.AssignableToTypeOf(time.Time{}), maxPerProcess, gomock.Any()).
					Return(populateMarvelCharacterResponse(totalData), nil).Times(1),
				store.
					EXPECT().
					Save(ctx, gomock.AssignableToTypeOf([]character.Character{})).
					Return(expectedErr),
			)

			err := m.InitializeMarvelCharacter(ctx)
			assert.Error(t, err)
			assert.True(t, errors.Cause(err) == expectedErr)
		})
	})

	//t.Run("When successfully save characters to database", func(t *testing.T) {
	//	const totalData = 10
	//
	//	resp := client.CharacterResponse{
	//		Data: client.Data{
	//			Total: totalData,
	//		},
	//	}
	//
	//	gomock.InOrder(
	//		store.
	//			EXPECT().
	//			CleanUp(ctx).
	//			Return(nil),
	//		marvelClient.
	//			EXPECT().
	//			Characters(ctx, gomock.AssignableToTypeOf(time.Time{}),1, 0).
	//			Return(resp, nil),
	//		marvelClient.
	//			EXPECT().
	//			Characters(ctx, gomock.AssignableToTypeOf(time.Time{}), maxPerProcess, gomock.Any()).
	//			Return(populateMarvelCharacterResponse(totalData), nil).Times(1),
	//		store.
	//			EXPECT().
	//			Save(ctx, gomock.AssignableToTypeOf([]character.Character{})).
	//			Return(nil),
	//	)
	//
	//	err := m.InitializeMarvelCharacter(ctx)
	//	assert.NoError(t, err)
	//})
}

func TestGetTotalCharacter(t *testing.T) {
	mockctl := gomock.NewController(t)
	defer mockctl.Finish()

	ctx := context.Background()
	marvelClient := mock.NewMockMarvelClient(mockctl)

	t.Run("When error occurred while getting characters from marvel API", func(t *testing.T) {
		t.Run("It should return error with 0 total data", func(t *testing.T) {
			expectedErr := errors.New("")
			marvelClient.
				EXPECT().
				Characters(ctx, gomock.AssignableToTypeOf(time.Time{}), 1, 0).
				Return(client.CharacterResponse{}, expectedErr)
			m := &marvelCmdExecImpl{
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
				Characters(ctx, gomock.AssignableToTypeOf(time.Time{}),1, 0).
				Return(resp, nil)

			m := &marvelCmdExecImpl{
				MarvelClient: marvelClient,
			}

			total, err := m.GetTotalCharacter(ctx)
			assert.NoError(t, err)
			assert.Equal(t, 10, total)
		})
	})
}
