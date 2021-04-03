package usecase

import (
	"context"
	"github.com/didiyudha/marvel/business/usecase/tests"
	"github.com/stretchr/testify/assert"
	"testing"

	teststore "github.com/didiyudha/marvel/business/data/tests"
)

var dbc = teststore.Container{
	Image: "postgres:13-alpine",
	Port:  "5432",
	Args:  []string{"-e", "POSTGRES_PASSWORD=postgres"},
}

var redisContainer = teststore.Container{
	Image: "redis",
	Port:  "6379",
}

func TestMarvel(t *testing.T) {
	store, caching, teardown := tests.NewStore(t, dbc, redisContainer)
	t.Cleanup(teardown)

	ctx := context.TODO()

	/**
	Note: There's seeding process within new store function above.
	 */

	marvelUseCase := NewMarvelUseCase(store, caching)

	characters, err := marvelUseCase.GetAllCharacterID(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, characters)
	assert.Equal(t, 3, len(characters))

	characterID := 1
	char, err := marvelUseCase.GetCharacter(ctx, characterID)
	assert.NoError(t, err)
	assert.Equal(t, characterID, char.ID)
}