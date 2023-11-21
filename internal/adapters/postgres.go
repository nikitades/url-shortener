package adapters

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/uptrace/opentelemetry-go-extra/otelsqlx"
)

func NewPgsqlConn(sqlconnstr string) (*sqlx.DB, error) {
	conn, err := otelsqlx.Connect("postgres", sqlconnstr)

	if err != nil {
		return nil, err
	}

	err = conn.Ping()

	if err != nil {
		return nil, err
	}

	return conn, nil
}

func Migrate(sqlconnstr string) error {
	conn, err := otelsqlx.Connect("postgres", sqlconnstr)
	if err != nil {
		return err
	}
	driver, err := postgres.WithInstance(conn.DB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/adapters/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}

	m.Up()

	return nil
}
