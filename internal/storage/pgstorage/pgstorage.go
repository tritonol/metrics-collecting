package pgstorage

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	m "github.com/tritonol/metrics-collecting.git/internal/structs/JSON"
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

func (pg *Postgres) CreateMetricTable(ctx context.Context) error {
	_, err := pg.db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS metrics (
			name varchar(128) not null,
			type varchar(32) not null,
			delta double precision,
			value integer,
			primary key(name,type)
		);
	`)

	if err != nil {
		return fmt.Errorf("cant create table: %w", err)
	}

	return nil
}

func (pg *Postgres) StoreMetric(ctx context.Context, name string, mType string, value float64, delta int64) error {
	_, err := pg.db.Exec(ctx, `
		INSERT INTO metrics(name, type, delta, value)
		VALUES($1, $2, $3, $4)
		ON CONFLICT (name,type)
		DO
			UPDATE SET delta = metrics.delta + $3, value = $4 WHERE metrics.name = $1 AND metrics.type = $2
	`, name, mType, delta, value)

	if err != nil {
		return err
	}

	return nil
}

func (pg *Postgres) GetMetric(ctx context.Context, name string, mType string) (m.Metric, error) {
	metric := m.Metric{}

	row := pg.db.QueryRow(ctx, `
		SELECT delta, value, type FROM metrics WHERE name = $1 AND type = $2  
	`, name, mType)

	err := row.Scan(&metric.Delta, &metric.Value, &metric.Type)
	if err != nil {
		return m.Metric{}, err
	}

	return metric, nil
}

func (pg *Postgres) GetMetrics(ctx context.Context) (map[string]m.Metric, error) {
	metrics := make(map[string]m.Metric, 29)
	rows, err := pg.db.Query(ctx, `SELECT name, delta, value, type FROM metrics`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var m m.Metric
		var name string

		rows.Scan(&name, &m.Delta, &m.Value, &m.Type)
		metrics[name] = m
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return metrics, err
}

func (pg *Postgres) GetAllDataStructed(ctx context.Context) (map[string]m.Metrics, error) {
	rawMetrics, err := pg.GetMetrics(ctx)
	if err != nil {
		return nil, err
	}

	preparedMetrics := make(map[string]m.Metrics)

	for name, metric := range rawMetrics {
		newMetric := metric
		preparedMetrics[name] = m.Metrics{
			ID:    name,
			MType: string(metric.Type),
			Value: &newMetric.Value,
			Delta: &newMetric.Delta,
		}
	}

	return preparedMetrics, nil
}

func (pg *Postgres) SaveAllDataStructured(ctx context.Context, metrics map[string]m.Metrics) error {
	for name, metric := range metrics {
		err := pg.StoreMetric(ctx, name, metric.MType, *metric.Value, *metric.Delta)
		if err != nil {
			return err
		}
	}

	return nil
}
