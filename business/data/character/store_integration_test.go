package character

import (
	"context"
	"github.com/bxcodec/faker/v3"
	"github.com/didiyudha/marvel/business/data/tests"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

var dbc = tests.Container{
	Image: "postgres:13-alpine",
	Port:  "5432",
	Args:  []string{"-e", "POSTGRES_PASSWORD=postgres"},
}

func TestCharacter(t *testing.T) {
	_, db, teardown := tests.NewDBContainer(t, dbc)
	t.Cleanup(teardown)

	ctx := context.TODO()

	store := NewStore(db)

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

	t.Run("Save Character", func(t *testing.T) {
		err := store.Save(ctx, characters...)
		assert.NoError(t, err)
	})

	t.Run("Find By ID", func(t *testing.T) {
		t.Run("When data found", func(t *testing.T) {
			t.Run("Should return correct data", func(t *testing.T) {
				characterID := characters[0].ID
				character, err := store.FindByID(ctx, characterID)
				assert.NoError(t, err)
				assert.Equal(t, characters[0].ID, character.ID)
				assert.Equal(t, characters[0].Name, character.Name)
				assert.Equal(t, characters[0].Description, character.Description)
			})
		})
		t.Run("When data not found", func(t *testing.T) {
			t.Run("Should return empty data with not found error", func(t *testing.T) {
				characterID := 100
				character, err := store.FindByID(ctx, characterID)
				assert.True(t, errors.Cause(err) == ErrNotFound)
				assert.Equal(t, 0, character.ID)
				assert.Equal(t, "", character.Name)
				assert.Equal(t, "", character.Description)
			})
		})
	})

	t.Run("Find All", func(t *testing.T) {
		t.Run("Should return all character data", func(t *testing.T) {
			characters, err := store.FindAll(ctx)
			assert.NoError(t, err)
			assert.True(t, len(characters) == 5)
		})
	})

	t.Run("Clean Up", func(t *testing.T) {
		err := store.CleanUp(ctx)
		assert.NoError(t, err)

		characters, err := store.FindAll(ctx)
		assert.NoError(t, err)
		assert.True(t, len(characters) == 0)
	})
}