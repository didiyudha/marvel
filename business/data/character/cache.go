package character

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"strconv"
	"time"
)



const (
	charactersKey = "characters"
)

type Caching interface {
	SetAll(ctx context.Context, characters []Character) error
	SetOne(ctx context.Context, character Character) error
	GetAll(ctx context.Context) ([]Character, error)
	FindOne(ctx context.Context, id int) (Character, error)
}

type cachingImpl struct {
	redisClient *redis.Client
}

func NewCaching(redisClient *redis.Client) Caching {
	return &cachingImpl{redisClient: redisClient}
}

func (c *cachingImpl) GetAll(ctx context.Context) ([]Character, error) {

	strRes, err := c.redisClient.Get(ctx, charactersKey).Result()

	if err == redis.Nil {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, errors.Wrap(err, "get all characters from caching")
	}

	var characters []Character

	if err := json.Unmarshal([]byte(strRes), &characters); err != nil {
		return nil, errors.Wrap(err, "unmarshal characters data from caching")
	}

	return characters, nil
}

func (c *cachingImpl) FindOne(ctx context.Context, id int) (Character, error) {

	strRes, err := c.redisClient.Get(ctx, strconv.Itoa(id)).Result()

	if err == redis.Nil {
		return Character{}, ErrNotFound
	}
	if err != nil {
		return Character{}, errors.Wrap(err, fmt.Sprintf("get character by id: %d", id))
	}

	var character Character

	if err := json.Unmarshal([]byte(strRes), &character); err != nil {
		return Character{}, errors.Wrap(err, "unmarshal single character data from redis")
	}

	return character, nil
}

func (c *cachingImpl) SetAll(ctx context.Context, characters []Character) error {

	b, err := json.Marshal(characters)
	if err != nil {
		return errors.Wrap(err, "unmarshal characters data")
	}

	return c.redisClient.Set(ctx, charactersKey, string(b), 5*time.Minute).Err()
}

func (c *cachingImpl) SetOne(ctx context.Context, character Character) error {

	b, err := json.Marshal(character)
	if err != nil {
		return errors.Wrap(err, "marshal character data")
	}

	return c.redisClient.Set(ctx, strconv.Itoa(character.ID), string(b), 5*time.Minute).Err()
}