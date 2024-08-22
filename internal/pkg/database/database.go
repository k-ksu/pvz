package database

import (
	"context"
	"fmt"
	"strconv"

	"HomeWork_1/internal/config"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Database struct {
	Cluster *pgxpool.Pool
}

func NewDatabase(ctx context.Context) (*Database, error) {
	pool, err := pgxpool.Connect(ctx, generateDsn())
	if err != nil {
		return nil, err
	}
	return &Database{Cluster: pool}, nil
}

func generateDsn() string {
	db, err := config.Read()
	if err != nil {
		return ""
	}

	port, _ := strconv.ParseInt(db.Port, 10, 64)
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", db.Host, port, db.User, db.Password, db.Name)
}

func (db Database) GetPool(_ context.Context) *pgxpool.Pool {
	return db.Cluster
}

func (db Database) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return pgxscan.Get(ctx, db.Cluster, dest, query, args...)
}

func (db Database) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return pgxscan.Select(ctx, db.Cluster, dest, query, args...)
}

func (db Database) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	return db.Cluster.Exec(ctx, query, args...)
}

func (db Database) ExecQueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return db.Cluster.QueryRow(ctx, query, args...)
}

func (db Database) BeginFunc(ctx context.Context, f func(pgx.Tx) error) error {
	return db.Cluster.BeginFunc(ctx, f)
}
