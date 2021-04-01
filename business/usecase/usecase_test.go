package usecase

import (
	"context"
	"github.com/bxcodec/faker/v3"
	"github.com/didiyudha/marvel/business/data/character"
	"github.com/didiyudha/marvel/business/data/character/mock"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMarvelUseCase(t *testing.T) {

	mockctl := gomock.NewController(t)
	defer mockctl.Finish()

	store := mock.NewMockStore(mockctl)
	caching := mock.NewMockCaching(mockctl)

	marvelUseCase := NewMarvelUseCase(store, caching)
	assert.NotNil(t, marvelUseCase)

	marvelImpl, ok := marvelUseCase.(*marvelUseCaseImpl)
	assert.True(t, ok)
	assert.NotNil(t, marvelImpl)
	assert.NotNil(t, marvelImpl.Store)
	assert.NotNil(t, marvelImpl.Cache)
}

func TestCollectCharacterID(t *testing.T) {
	characters := []character.Character{
		{
			ID:          1,
			Name:        faker.Name(),
			Description: faker.Paragraph(),
		},
		{
			ID:          2,
			Name:        faker.Name(),
			Description: faker.Paragraph(),
		},
		{
			ID:          3,
			Name:        faker.Name(),
			Description: faker.Paragraph(),
		},
		{
			ID:          4,
			Name:        faker.Name(),
			Description: faker.Paragraph(),
		},
		{
			ID:          5,
			Name:        faker.Name(),
			Description: faker.Paragraph(),
		},
	}
	marvelImpl := new(marvelUseCaseImpl)
	ids := marvelImpl.collectCharacterID(characters)
	expectedID := []int{1, 2, 3, 4, 5}

	assert.Equal(t, expectedID, ids)

	characters = []character.Character{}
	ids = marvelImpl.collectCharacterID(characters)
	assert.Equal(t, []int{}, ids)

}

func TestGetAllCharacterID(t *testing.T) {
	mockctl := gomock.NewController(t)
	defer mockctl.Finish()

	characters := []character.Character{
		{
			ID:          1,
			Name:        faker.Name(),
			Description: faker.Paragraph(),
		},
		{
			ID:          2,
			Name:        faker.Name(),
			Description: faker.Paragraph(),
		},
		{
			ID:          3,
			Name:        faker.Name(),
			Description: faker.Paragraph(),
		},
		{
			ID:          4,
			Name:        faker.Name(),
			Description: faker.Paragraph(),
		},
		{
			ID:          5,
			Name:        faker.Name(),
			Description: faker.Paragraph(),
		},
	}

	t.Run("Successfully got characters from caching", func(t *testing.T) {
		ctx := context.Background()
		caching := mock.NewMockCaching(mockctl)
		caching.
			EXPECT().
			GetAll(ctx).
			Return(characters, nil)
		m := &marvelUseCaseImpl{
			Cache: caching,
		}
		ids, err := m.GetAllCharacterID(ctx)
		assert.NoError(t, err)
		assert.Equal(t, []int{1, 2, 3, 4, 5}, ids)
	})

	t.Run("Error occurred when getting data from database", func(t *testing.T) {
		ctx := context.Background()
		caching := mock.NewMockCaching(mockctl)
		store := mock.NewMockStore(mockctl)
		expectedErr := errors.New("database down")

		gomock.InOrder(
			caching.
				EXPECT().
				GetAll(ctx).
				Return(nil, nil),
			store.
				EXPECT().
				FindAll(ctx).
				Return(nil, expectedErr),
		)

		m := &marvelUseCaseImpl{
			Cache: caching,
			Store: store,
		}

		ids, err := m.GetAllCharacterID(ctx)
		assert.Error(t, err)
		assert.True(t, errors.Cause(err) == expectedErr)
		assert.Nil(t, ids)
	})

	t.Run("Successfully got data from database", func(t *testing.T) {
		ctx := context.Background()
		caching := mock.NewMockCaching(mockctl)
		store := mock.NewMockStore(mockctl)

		gomock.InOrder(
			caching.
				EXPECT().
				GetAll(ctx).
				Return(nil, nil),
			store.
				EXPECT().
				FindAll(ctx).
				Return(characters, nil),
			caching.
				EXPECT().
				SetAll(ctx, gomock.AssignableToTypeOf([]character.Character{})),
		)

		m := &marvelUseCaseImpl{
			Cache: caching,
			Store: store,
		}
		ids, err := m.GetAllCharacterID(ctx)
		assert.NoError(t, err)
		assert.Equal(t, []int{1, 2, 3, 4, 5}, ids)
	})
}

func TestGetCharacter(t *testing.T) {
	mockctl := gomock.NewController(t)
	defer mockctl.Finish()

	ctx := context.Background()
	id := 1

	char := character.Character{
		ID:          id,
		Name:        faker.Name(),
		Description: faker.Paragraph(),
	}

	t.Run("Successfully got the character from caching", func(t *testing.T) {
		caching := mock.NewMockCaching(mockctl)
		caching.
			EXPECT().
			FindOne(ctx, id).
			Return(char, nil)

		m := &marvelUseCaseImpl{
			Cache: caching,
		}

		charResp, err := m.GetCharacter(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, char.ID, charResp.ID)
		assert.Equal(t, char.Name, charResp.Name)
		assert.Equal(t, char.Description, charResp.Description)
	})

	t.Run("", func(t *testing.T) {
		store := mock.NewMockStore(mockctl)
		caching := mock.NewMockCaching(mockctl)

		gomock.InOrder(
			caching.
				EXPECT().
				FindOne(ctx, id).
				Return(character.Character{}, nil),
			store.
				EXPECT().
				FindByID(ctx, id).
				Return(character.Character{}, character.ErrNotFound),
		)

		m := &marvelUseCaseImpl{
			Store: store,
			Cache: caching,
		}

		charResp, err := m.GetCharacter(ctx, id)
		assert.Equal(t, character.ErrNotFound, err)
		assert.Equal(t, 0, charResp.ID)
		assert.Equal(t, "", charResp.Name)
		assert.Equal(t, "", charResp.Description)
	})

	t.Run("Error occurred when find character from database", func(t *testing.T) {
		store := mock.NewMockStore(mockctl)
		caching := mock.NewMockCaching(mockctl)
		expectedErr := errors.New("database down")

		gomock.InOrder(
			caching.
				EXPECT().
				FindOne(ctx, id).
				Return(character.Character{}, nil),
			store.
				EXPECT().
				FindByID(ctx, id).
				Return(character.Character{}, expectedErr),
		)

		m := &marvelUseCaseImpl{
			Store: store,
			Cache: caching,
		}

		charResp, err := m.GetCharacter(ctx, id)
		assert.Error(t, err)
		assert.Equal(t, 0, charResp.ID)
		assert.Equal(t, "", charResp.Name)
		assert.Equal(t, "", charResp.Description)
	})

	t.Run("Successfully get character from database", func(t *testing.T) {
		store := mock.NewMockStore(mockctl)
		caching := mock.NewMockCaching(mockctl)

		gomock.InOrder(
			caching.
				EXPECT().
				FindOne(ctx, id).
				Return(character.Character{}, nil),
			store.
				EXPECT().
				FindByID(ctx, id).
				Return(char, nil),
			caching.
				EXPECT().
				SetOne(ctx, gomock.AssignableToTypeOf(character.Character{})),
		)

		m := &marvelUseCaseImpl{
			Store: store,
			Cache: caching,
		}

		charResp, err := m.GetCharacter(ctx, id)

		assert.NoError(t, err)
		assert.Equal(t, char.ID, charResp.ID)
		assert.Equal(t, char.Name, charResp.Name)
		assert.Equal(t, char.Description, charResp.Description)
	})
}
