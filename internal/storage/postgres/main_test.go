package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"

	"HomeWork_1/internal/pkg/database"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/pressly/goose/v3"
)

var (
	Host     = "localhost"
	Name     = "test"
	Password = "test"
	Port     = "5432"
	User     = "test"

	migrationDir = "../../../migrations"

	testRepo *Repository
)

func NewTestDocker() (*database.Database, func()) {
	var (
		dsn      string
		resource *dockertest.Resource
		db       *pgxpool.Pool
		ctx      = context.Background()
	)

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatal("cannot connect to test docker ", err)
	}

	resource, err = pool.Run("postgres", "13", []string{
		fmt.Sprintf("POSTGRES_USER=%s", User),
		fmt.Sprintf("POSTGRES_PASSWORD=%s", Password),
		fmt.Sprintf("POSTGRES_DB=%s", Name),
	})
	if err != nil {
		log.Fatal("pool run failed", err)
	}

	closeFunc := func() {
		if err := pool.Purge(resource); err != nil {
			log.Fatal("cannot close resource", err)
		}
	}

	PortDB, _ := strconv.Atoi(resource.GetPort("5432/tcp"))
	err = pool.Retry(func() error {
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", Host, PortDB, User, Password, Name)
		db, err = pgxpool.Connect(ctx, dsn)
		if err != nil {
			return errors.New("testDB connection failed: " + err.Error())
		}

		err = db.Ping(ctx)
		if err != nil {
			return errors.New("ping failed: " + err.Error())
		}

		dbSql, err := sql.Open("postgres", dsn)
		defer dbSql.Close()
		if err != nil {
			return errors.New("dbSql failed: " + err.Error())
		}
		goose.SetTableName("goose_db_version")
		err = goose.SetDialect("postgres")
		if err != nil {
			return errors.New("set dialect failed: " + err.Error())
		}
		err = goose.Up(dbSql, migrationDir)
		if err != nil {
			return errors.New("migration up failed: " + err.Error())
		}

		return nil
	})

	if err != nil {
		closeFunc()
		log.Fatal(err)
	}

	return &database.Database{Cluster: db}, closeFunc
}

func TestMain(m *testing.M) {
	var exitCode int
	func() {
		db, closeFunc := NewTestDocker()
		defer closeFunc()

		testRepo = NewRepository(db)
		exitCode = m.Run()
	}()

	os.Exit(exitCode)
}
