package usecase

import (
	"context"
	"github.com/didiyudha/marvel/business/data/character"
	"github.com/pkg/errors"
	"log"
)

type MarvelUseCase interface {
	GetAllCharacterID(ctx context.Context) ([]int, error)
	GetCharacter(ctx context.Context, id int) (character.Character, error)
}

type marvelUseCaseImpl struct {
	Cache character.Caching
	Store character.Store
}

func NewMarvelUseCase(store character.Store, caching character.Caching) MarvelUseCase {
	return &marvelUseCaseImpl{
		Store: store,
		Cache: caching,
	}
}

func (m *marvelUseCaseImpl) collectCharacterID(characters []character.Character) (id []int) {
	if len(characters) == 0 {
		return []int{}
	}
	for _, c := range characters {
		id = append(id, c.ID)
	}
	return
}

func (m *marvelUseCaseImpl) GetAllCharacterID(ctx context.Context) ([]int, error) {

	characters, err := m.Cache.GetAll(ctx)

	if err != nil && err != character.ErrNotFound {
		// Silent error and log the error.
		log.Printf("[ERROR] error when getting data all characters from cache")
	}

	if characters != nil && len(characters) > 0 {
		return m.collectCharacterID(characters), nil
	}

	characters, err = m.Store.FindAll(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "find all characters")
	}

	if err := m.Cache.SetAll(ctx, characters); err != nil {
		// Silent error and log the error.
		log.Printf("[ERROR] error when set all characters to chaching")
	}

	return m.collectCharacterID(characters), nil
}

func (m *marvelUseCaseImpl) GetCharacter(ctx context.Context, id int) (character.Character, error) {

	char, err := m.Cache.FindOne(ctx, id)
	if err != nil && err != character.ErrNotFound {
		// Silent error and log the error.
		log.Printf("[ERROR] error when getting data a character from cache")
	}

	if char.ID > 0 && char.Name != "" {
		return char, nil
	}

	char, err = m.Store.FindByID(ctx, id)
	if err == character.ErrNotFound {
		return character.Character{}, err
	}
	if err != nil {
		return character.Character{}, errors.Wrap(err, "find character by id")
	}

	if err := m.Cache.SetOne(ctx, char); err != nil {
		log.Printf("[ERROR] error when set one character to chache")
	}

	return char, nil
}
