package caching

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"time"
)

type Config struct {
	Addr string
	Password string
	DB int
}

func New(cfg Config) (*redis.Client, error) {

	client := redis.NewClient(&redis.Options{
		Addr: cfg.Addr,
		Password: cfg.Password,
		DB: cfg.DB,
	})

	err := client.Ping(context.Background()).Err()

	threshold := 5
	counter := 1

	for  counter <= threshold {
		client = redis.NewClient(&redis.Options{
			Addr: cfg.Addr,
			Password: cfg.Password,
			DB: cfg.DB,
		})
		err = client.Ping(context.Background()).Err()
		if err == nil {
			break
		}
		time.Sleep(time.Duration(counter)*100*time.Millisecond)
	}

	if err != nil {
		return nil, errors.Wrap(err, "unable to connect to redis server")
	}

	return client, nil
}
