package pgstorage

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	db *pgxpool.Pool
}

var (
	pgInstance *Postgres
	pgOnce     sync.Once
)

func NewPg(ctx context.Context, connString string) (*Postgres, error) {
	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, connString)
		if err != nil {
			return
		}

		pgInstance = &Postgres{db}
	})

	return pgInstance, nil
}

func (pg *Postgres) Ping(ctx context.Context) error {
	return pg.db.Ping(ctx)
}

func (pg *Postgres) Close() {
	pg.db.Close()
}

func (pg *Postgres) CreateMetricTable(ctx context.Context) error{
	_, err := pg.db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS metrics (
			id serial primary key,
			name varchar(128) not null,
			type varchar(32) not null,
			delta double precision,
			value integer
		);
	`)

	if err != nil {
		return fmt.Errorf("cant create table: %w", err)
	}

	return nil
}
