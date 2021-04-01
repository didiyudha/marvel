package internal

import (
	"context"
	"fmt"
	"github.com/didiyudha/marvel/business/data/character"
	"github.com/didiyudha/marvel/client"
	"github.com/pkg/errors"
	"log"
	"math"
	"time"
)

const (
	maxPerProcess = 100
)

type MarvelWorker interface {
	Do()
}

type marvelWorkerImpl struct {
	Store character.Store
	MarvelClient client.MarvelClient
}

func NewMarvelWorker(store character.Store, marvelClient client.MarvelClient) MarvelWorker {
	return &marvelWorkerImpl{
		Store: store,
		MarvelClient: marvelClient,
	}
}

func findNewData(fromAPI []character.Character, fromDB []character.Character) []character.Character {

	idMapDB := make(map[int]bool)
	for _, c := range fromDB {
		idMapDB[c.ID] = true
	}

	newCharacters := make([]character.Character, 0, len(fromAPI))
	for _, char := range fromAPI {
		_, ok := idMapDB[char.ID]
		if !ok {
			newCharacters = append(newCharacters, char)
		}
	}

	return newCharacters
}

func (m *marvelWorkerImpl) Do() {
	fmt.Println("worker started")
	ctx := context.Background()

	charactersFromAPI, err := m.GetAllCharactersFromAPI(ctx)
	if err != nil {
		log.Printf("[ERROR] gt all characters from API: %v\n", err)
		return
	}

	charactersFromDB, err := m.Store.FindAll(context.Background())
	if err != nil {
		log.Printf("[ERROR] gt all characters from database: %v\n", err)
	}

	newCharacters := findNewData(charactersFromAPI, charactersFromDB)

	if len(newCharacters) > 0 {
		if err := m.Store.Save(ctx, newCharacters...); err != nil {
			fmt.Printf("[ERROR] save new characters by worker: %v\n", err)
			return
		}
	}
	fmt.Println("worker started")
}

func (m *marvelWorkerImpl) GetTotalCharacter(ctx context.Context) (total int, err error) {
	ts := time.Now()
	res, err := m.MarvelClient.Characters(ctx, ts, 1, 0)

	if err != nil {
		err = errors.Wrap(err, "get total characters")
		return
	}

	total = res.Data.Total
	return
}

func (m *marvelWorkerImpl) GetAllCharactersFromAPI(ctx context.Context) ([]character.Character, error) {

	total, err := m.GetTotalCharacter(ctx)
	if err != nil {
		return nil, err
	}

	totalProcessing := int(math.Ceil(float64(total) / float64(maxPerProcess)))

	var offset int
	var results []client.Result

	for i := 1; i <= totalProcessing; i++ {
		ts := time.Now()
		resp, err := m.MarvelClient.Characters(ctx, ts, maxPerProcess, offset)
		if err != nil {
			return nil, errors.Wrap(err, "get characters")
		}
		results = append(results, resp.Data.Results...)
		offset = offset + maxPerProcess
	}

	characters := make([]character.Character, 0, len(results))

	for _, r := range results {
		c := character.Character{
			ID:          r.ID,
			Name:        r.Name,
			Description: r.Description,
		}
		characters = append(characters, c)
	}

	return characters, nil

}
