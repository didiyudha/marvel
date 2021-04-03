package tests

import (
	"context"
	"github.com/didiyudha/marvel/business/data/character"
	"github.com/didiyudha/marvel/business/data/schema"
	"github.com/didiyudha/marvel/business/data/tests"
	"testing"
)

func NewStore(t *testing.T, dbContainer tests.Container, redisContainer tests.Container) (character.Store, character.Caching, func()) {

	_, db, dbTeardown := tests.NewDBContainer(t, dbContainer)
	_, redis, redisTeardown := tests.NewCachingContainer(t, redisContainer)

	store := character.NewStore(db)
	caching := character.NewCaching(redis)

	if err := schema.Seed(context.Background(), db); err != nil {
		t.Fatal(err)
	}

	teardown := func() {
		t.Helper()
		schema.DeleteAll(context.Background(), db)
		dbTeardown()
		redisTeardown()
	}

	return store, caching, teardown
}