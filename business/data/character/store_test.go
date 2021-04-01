package character

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bxcodec/faker/v3"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSave(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	ctx := context.Background()

	char := Character{
		ID:          1,
		Name:        faker.Name(),
		Description: faker.Paragraph(),
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	t.Run("Successfully save character to database", func(t *testing.T) {
		store := NewStore(sqlxDB)
		result := sqlmock.NewResult(1, 1)
		mock.
			ExpectExec("INSERT INTO characters").
			WillReturnResult(result)

		err = store.Save(ctx, char)
		assert.NoError(t, err)
	})
}

func TestFindByID(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	ctx := context.Background()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	characterID := 1

	char := Character{
		ID:          characterID,
		Name:        faker.Name(),
		Description: faker.Paragraph(),
	}

	t.Run("When data not found", func(t *testing.T) {
		t.Run("It should return empty character data with error not found signature", func(t *testing.T) {
			store := NewStore(sqlxDB)

			mock.
				ExpectQuery("^SELECT (.+) FROM characters").
				WithArgs(characterID).
				WillReturnError(sql.ErrNoRows)

			char, err := store.FindByID(ctx, characterID)
			assert.Error(t, err)
			assert.Equal(t, 0, char.ID)
			assert.Equal(t, "", char.Name)
			assert.Equal(t, "", char.Description)
		})
	})

	t.Run("When successfully get character data from database", func(t *testing.T) {
		t.Run("It should successfully return a character data", func(t *testing.T) {
			store := NewStore(sqlxDB)
			rows := sqlmock.NewRows([]string{"id", "name", "description"}).AddRow(char.ID, char.Name, char.Description)
			mock.
				ExpectQuery("^SELECT (.+) FROM characters").
				WithArgs(characterID).
				WillReturnRows(rows).
				WillReturnError(nil)

			charResp, err := store.FindByID(ctx, characterID)
			assert.NoError(t, err)
			assert.Equal(t, char.ID, charResp.ID)
			assert.Equal(t, char.Name, charResp.Name)
			assert.Equal(t, char.Description, charResp.Description)
		})
	})
}

func TestFindAll(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	ctx := context.Background()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	t.Run("When error occurred while getting all characters", func(t *testing.T) {
		t.Run("It should return error", func(t *testing.T) {
			store := NewStore(sqlxDB)
			expectedErr := errors.New("database down")
			mock.
				ExpectQuery("^SELECT (.+) FROM characters").
				WillReturnError(expectedErr)

			chars, err := store.FindAll(ctx)
			assert.True(t, errors.Cause(err) == expectedErr)
			assert.Nil(t, chars)
		})
	})

	t.Run("When successfully get characters from database", func(t *testing.T) {
		t.Run("It should return all characters data", func(t *testing.T) {
			store := NewStore(sqlxDB)
			rows := sqlmock.NewRows([]string{"id", "name", "description"}).
				AddRow(1, faker.Name(), faker.Paragraph()).
				AddRow(2, faker.Name(), faker.Paragraph()).
				AddRow(3, faker.Name(), faker.Paragraph()).
				AddRow(4, faker.Name(), faker.Paragraph()).
				AddRow(5, faker.Name(), faker.Paragraph())
			mock.
				ExpectQuery("^SELECT (.+) FROM characters").
				WillReturnRows(rows).
				WillReturnError(nil)

			chars, err := store.FindAll(ctx)
			assert.NoError(t, err)
			assert.Equal(t, 5, len(chars))
		})
	})
}

func TestCleanUp(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	ctx := context.Background()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	store := NewStore(sqlxDB)

	t.Run("When error clean up data", func(t *testing.T) {
		t.Run("It should return error", func(t *testing.T) {
			expectedErr := errors.New("database down")
			mock.
				ExpectExec(`DELETE FROM characters`).
				WillReturnResult(sqlmock.NewResult(0, 0)).
				WillReturnError(expectedErr)
			err = store.CleanUp(ctx)
			assert.Error(t, err)
			assert.True(t, errors.Cause(err) == expectedErr)
		})
	})

	t.Run("When successfully clean up the data", func(t *testing.T) {
		t.Run("It should return no error", func(t *testing.T) {
			mock.
				ExpectExec(`DELETE FROM characters`).
				WillReturnResult(sqlmock.NewResult(0, 0)).
				WillReturnError(nil)
			err = store.CleanUp(ctx)
			assert.NoError(t, err)
		})
	})
}