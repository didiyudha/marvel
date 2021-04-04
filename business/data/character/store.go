package character

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var (
	ErrNotFound = errors.New("not found")
)

type Store interface {
	CleanUp(ctx context.Context) error
	Save(ctx context.Context, newCharacter ...Character) error
	FindByID(ctx context.Context, id int) (Character, error)
	FindAll(ctx context.Context) ([]Character, error)
}

type persistentStorage struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) Store {
	return &persistentStorage{db: db}
}

func (p *persistentStorage) Save(ctx context.Context, newCharacter ...Character) error {

	q := `INSERT INTO characters (id, name, description) VALUES `

	insertParams := make([]interface{}, 0, len(newCharacter))

	for i, v := range newCharacter {
		p1 := i * 3
		q += fmt.Sprintf("($%d,$%d,$%d),", p1+1,p1+2,p1+3)
		insertParams = append(insertParams, v.ID, v.Name, v.Description)
	}

	q = q[:len(q)-1]

	_, err := p.db.MustExecContext(ctx, q, insertParams...).RowsAffected()
	if err != nil {
		return errors.Wrap(err, "save new character")
	}

	return nil
}

func (p *persistentStorage) FindByID(ctx context.Context, id int) (Character, error) {

	var character Character

	q := `SELECT id,
			name,
			description
		FROM characters 
		WHERE id = $1`

	err := p.db.GetContext(ctx, &character, q, id)
	if err == sql.ErrNoRows {
		return Character{}, ErrNotFound
	}
	if err != nil {
		return Character{}, errors.Wrap(err, "find character by id")
	}

	return character, nil
}

func (p *persistentStorage) FindAll(ctx context.Context) ([]Character, error) {

	var characters []Character

	q := `SELECT id,
			name,
			description
		FROM "characters"`

	if err := p.db.SelectContext(ctx, &characters, q); err != nil {
		return nil, errors.Wrap(err, "select all characters")
	}

	return characters, nil
}

func (p *persistentStorage) CleanUp(ctx context.Context) error {
	q := `DELETE FROM characters`
	_, err := p.db.ExecContext(ctx, q)
	if err != nil {
		return errors.Wrap(err, "execute clean up query")
	}
	return nil
}