package character

import (
	"context"
	"github.com/bxcodec/faker/v3"
	"github.com/didiyudha/marvel/business/data/tests"
	"github.com/stretchr/testify/assert"
	"testing"
)

var redisContainer = tests.Container{
	Image: "redis",
	Port:  "6379",
}

func TestCaching(t *testing.T) {
	_, redisClient, teardown := tests.NewCachingContainer(t, redisContainer)
	t.Cleanup(teardown)

	ctx := context.TODO()
	cache := NewCaching(redisClient)

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

	result, err := cache.GetAll(ctx)
	assert.Equal(t, ErrNotFound, err)
	assert.Nil(t, result)

	err = cache.SetAll(ctx, characters)
	assert.NoError(t, err)

	result, err = cache.GetAll(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 5, len(result))

	assert.Equal(t, characters[0].ID, result[0].ID)
	assert.Equal(t, characters[0].Name, result[0].Name)
	assert.Equal(t, characters[0].Description, result[0].Description)

	assert.Equal(t, characters[1].ID, result[1].ID)
	assert.Equal(t, characters[1].Name, result[1].Name)
	assert.Equal(t, characters[1].Description, result[1].Description)

	assert.Equal(t, characters[2].ID, result[2].ID)
	assert.Equal(t, characters[2].Name, result[2].Name)
	assert.Equal(t, characters[2].Description, result[2].Description)

	assert.Equal(t, characters[3].ID, result[3].ID)
	assert.Equal(t, characters[3].Name, result[3].Name)
	assert.Equal(t, characters[3].Description, result[3].Description)

	assert.Equal(t, characters[4].ID, result[4].ID)
	assert.Equal(t, characters[4].Name, result[4].Name)
	assert.Equal(t, characters[4].Description, result[4].Description)


	char := characters[0]
	err = cache.SetOne(ctx, char)
	assert.NoError(t, err)

	characterID := characters[0].ID
	character, err := cache.FindOne(ctx, characterID)
	assert.NoError(t, err)
	assert.Equal(t, characters[0].ID, character.ID)

}
