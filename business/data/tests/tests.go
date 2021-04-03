package tests

import (
	"context"
	"github.com/didiyudha/marvel/business/data/schema"
	"github.com/didiyudha/marvel/foundation/database"
	"github.com/didiyudha/marvel/foundation/docker"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
	"testing"
	"time"
)

// Success and failure markers.
const (
	Success = "\u2713"
	Failed  = "\u2717"
)

type DBContainer struct {
	Image string
	Port string
	Args []string
}

func NewDBContainer(t *testing.T, dbc DBContainer) (logger *log.Logger, db *sqlx.DB, teardown func()) {
	c := docker.StartContainer(t, dbc.Image, dbc.Port, dbc.Args...)

	var err error

	db, err = database.Open(database.Config{
		User:         "postgres",
		Password:     "postgres",
		Host:         c.Host,
		Name:         "postgres",
		DisableTLS:   true,
	})

	if err != nil {
		t.Fatalf("Opening database connection: %v", err)
	}

	t.Log("Waiting for database to be ready")

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	if err := schema.Migrate(ctx, db); err != nil {
		docker.DumpContainerLogs(t, c.ID)
		docker.StopContainer(t, c.ID)
		t.Fatalf("Migrating error: %s", c.ID)
	}

	teardown = func() {
		t.Helper()
		db.Close()
		docker.StopContainer(t, c.ID)
	}

	logger = log.New(os.Stdout, "TEST", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	return
}