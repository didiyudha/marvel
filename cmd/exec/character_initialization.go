package exec

import (
	"context"
	"github.com/didiyudha/marvel/business/data/character"
	"github.com/didiyudha/marvel/client"
	"github.com/pkg/errors"
	"math"
	"time"
)

const (
	maxPerProcess = 100
)

type MarvelCmdExecutor interface {
	InitializeMarvelCharacter(ctx context.Context) error
}

type marvelCmdExecImpl struct {
	Store character.Store
	MarvelClient client.MarvelClient
}



func (m *marvelCmdExecImpl) InitializeMarvelCharacter(ctx context.Context) error {

	if err := m.Store.CleanUp(ctx); err != nil {
		return errors.Wrap(err, "clean up characters table")
	}

	total, err := m.GetTotalCharacter(ctx)
	if err != nil {
		return err
	}

	totalProcessing := int(math.Ceil(float64(total) / float64(maxPerProcess)))

	var offset int
	var results []client.Result

	for i := 1; i <= totalProcessing; i++ {
		ts := time.Now()
		resp, err := m.MarvelClient.Characters(ctx, ts, maxPerProcess, offset)
		if err != nil {
			return errors.Wrap(err, "get characters")
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

	if err := m.Store.Save(ctx, characters...); err != nil {
		return errors.Wrapf(err, "save characters")
	}
	return nil
}

func (m *marvelCmdExecImpl) GetTotalCharacter(ctx context.Context) (total int, err error) {
	ts := time.Now()
	res, err := m.MarvelClient.Characters(ctx, ts, 1, 0)

	if err != nil {
		err = errors.Wrap(err, "get total characters")
		return
	}

	total = res.Data.Total
	return
}