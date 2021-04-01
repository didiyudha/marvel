package exec

import (
	"context"
	"fmt"
	"github.com/didiyudha/marvel/business/data/schema"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"
)

type TableMigrator interface {
	MigrateTable(ctx context.Context) error
	DeleteAll(ctx context.Context) error
}

type tableMigratorImpl struct {
	DB *sqlx.DB
}

func NewTableMigrator(db *sqlx.DB) TableMigrator {
	return &tableMigratorImpl{
		DB: db,
	}
}

func (t *tableMigratorImpl) MigrateTable(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	if err := schema.Migrate(ctx, t.DB); err != nil {
		return errors.Wrap(err, "table migration")
	}

	fmt.Println("migration complete")
	return nil
}

func (t *tableMigratorImpl) DeleteAll(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	if err := schema.DeleteAll(ctx, t.DB); err != nil {
		return errors.Wrapf(err, "delete all characters")
	}

	fmt.Println("delete all characters complete")

	return nil
}