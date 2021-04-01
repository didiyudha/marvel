package character

import (
	"context"
	"encoding/json"
	"github.com/alicebob/miniredis"
	"github.com/bxcodec/faker/v3"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/pkg/errors"
	"strconv"
	"time"

	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCaching(t *testing.T) {
	mr, err := miniredis.Run()
	assert.NoError(t, err)
	defer mr.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	caching := NewCaching(redisClient)
	assert.NotNil(t, caching)

	impl, ok := (caching).(*cachingImpl)
	assert.True(t, ok)
	assert.NotNil(t, impl)
	assert.NotNil(t, impl.redisClient)
}

func TestGetAll(t *testing.T) {
	db, mock := redismock.NewClientMock()
	ctx := context.TODO()

	caching := NewCaching(db)

	t.Run("When key not found", func(t *testing.T) {
		t.Run("It should return error not found", func(t *testing.T) {
			mock.ExpectGet(charactersKey).RedisNil()
			characters, err := caching.GetAll(ctx)
			assert.Error(t, err)
			assert.True(t, errors.Cause(err) == ErrNotFound)
			assert.Nil(t, characters)
		})
	})

	t.Run("When successfully get the data", func(t *testing.T) {
		t.Run("It should return the data", func(t *testing.T) {
			characters := []Character{
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
			b, err := json.Marshal(characters)
			assert.NoError(t, err)

			mock.ExpectGet(charactersKey).SetVal(string(b))

			resp, err := caching.GetAll(ctx)
			assert.NoError(t, err)
			assert.Equal(t, 5, len(resp))

			assert.Equal(t, characters[0].ID, resp[0].ID)
			assert.Equal(t, characters[0].Name, resp[0].Name)
			assert.Equal(t, characters[0].Description, resp[0].Description)

			assert.Equal(t, characters[1].ID, resp[1].ID)
			assert.Equal(t, characters[1].Name, resp[1].Name)
			assert.Equal(t, characters[1].Description, resp[1].Description)

			assert.Equal(t, characters[2].ID, resp[2].ID)
			assert.Equal(t, characters[2].Name, resp[2].Name)
			assert.Equal(t, characters[2].Description, resp[2].Description)

			assert.Equal(t, characters[3].ID, resp[3].ID)
			assert.Equal(t, characters[3].Name, resp[3].Name)
			assert.Equal(t, characters[3].Description, resp[3].Description)

			assert.Equal(t, characters[4].ID, resp[4].ID)
			assert.Equal(t, characters[4].Name, resp[4].Name)
			assert.Equal(t, characters[4].Description, resp[4].Description)

		})
	})
}

func TestFindOne(t *testing.T) {
	db, mock := redismock.NewClientMock()
	ctx := context.TODO()

	caching := NewCaching(db)
	characterID := 1
	character := Character{
		ID:          1,
		Name:        faker.Name(),
		Description: faker.Paragraph(),
	}

	t.Run("When key not found", func(t *testing.T) {
		t.Run("It should return error not found", func(t *testing.T) {
			mock.
				ExpectGet(strconv.Itoa(characterID)).
				RedisNil()

			char, err := caching.FindOne(ctx, characterID)
			assert.Equal(t, err, ErrNotFound)
			assert.Equal(t, 0, char.ID)
			assert.Equal(t, "", char.Name)
			assert.Equal(t, "", char.Description)
		})
	})
	t.Run("When key found", func(t *testing.T) {
		t.Run("It should return character data", func(t *testing.T) {
			b, err := json.Marshal(character)
			assert.NoError(t, err)

			mock.
				ExpectGet(strconv.Itoa(characterID)).
				SetVal(string(b))

			char, err := caching.FindOne(ctx, characterID)
			assert.NoError(t, err)
			assert.Equal(t, character.ID, char.ID)
			assert.Equal(t, character.Name, char.Name)
			assert.Equal(t, character.Description, char.Description)

		})
	})
}

func TestSetAll(t *testing.T) {
	db, mock := redismock.NewClientMock()
	ctx := context.TODO()

	caching := NewCaching(db)

	characters := []Character{
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

	t.Run("When error occurred while set data to caching", func(t *testing.T) {

		t.Run("It should  return error", func(t *testing.T) {
			b, err := json.Marshal(characters)
			assert.NoError(t, err)

			expectedErr := errors.New("redis down")

			mock.
				ExpectSet(charactersKey, string(b), 5*time.Minute).
				SetErr(expectedErr)

			err = caching.SetAll(ctx, characters)
			assert.Equal(t, expectedErr, err)

		})

	})

	t.Run("When successfully get the data", func(t *testing.T) {
		t.Run("It should return no error", func(t *testing.T) {
			b, err := json.Marshal(characters)
			assert.NoError(t, err)

			mock.
				ExpectSet(charactersKey, string(b), 5*time.Minute).
				SetVal("")

			err = caching.SetAll(ctx, characters)
			assert.NoError(t, err)
		})
	})
}

func TestSetOne(t *testing.T) {
	db, mock := redismock.NewClientMock()
	ctx := context.TODO()

	caching := NewCaching(db)
	character := Character{
		ID:          1,
		Name:        faker.Name(),
		Description: faker.Paragraph(),
	}

	t.Run("When error occurred while set character data", func(t *testing.T) {
		t.Run("It should return error", func(t *testing.T) {
			b, err := json.Marshal(character)
			assert.NoError(t, err)
			expectedErr := errors.New("redis down")
			mock.
				ExpectSet(strconv.Itoa(character.ID), string(b), 5 * time.Minute).
				SetErr(expectedErr)

			err = caching.SetOne(ctx, character)
			assert.Error(t, err)
			assert.Equal(t, "redis down", err.Error())

		})
	})

	t.Run("When successfully set character data to caching", func(t *testing.T) {
		t.Run("It should return no error", func(t *testing.T) {
			b, err := json.Marshal(character)
			assert.NoError(t, err)

			mock.
				ExpectSet(strconv.Itoa(character.ID), string(b), 5 * time.Minute).
				SetVal("")

			err = caching.SetOne(ctx, character)
			assert.NoError(t, err)
		})
	})
}
