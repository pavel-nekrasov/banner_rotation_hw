package storage

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/stdlib" // need import pgx
	goose "github.com/pressly/goose/v3"
)

type Storage struct {
	dsn string
	db  *sql.DB
}

func New(host string, port int, dbname, user, password string) *Storage {
	return &Storage{
		dsn: fmt.Sprintf("postgres://%v:%v@%v:%v/%v", user, password, host, port, dbname),
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	var err error
	s.db, err = sql.Open("pgx", s.dsn)
	if err != nil {
		return fmt.Errorf("cannot open pgx driver: %w", err)
	}

	return s.db.PingContext(ctx)
}

func (s *Storage) Close(_ context.Context) error {
	return s.db.Close()
}

func (s *Storage) Migrate(_ context.Context, migrate string) (err error) {
	//	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("cannot set dialect: %w", err)
	}

	if err := goose.Up(s.db, migrate); err != nil {
		return fmt.Errorf("cannot do up migration: %w", err)
	}

	return nil
}
