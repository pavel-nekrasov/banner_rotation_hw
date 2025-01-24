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
	DB  *sql.DB
}

func New(host string, port int, dbname, user, password string) *Storage {
	return &Storage{
		dsn: fmt.Sprintf("postgres://%v:%v@%v:%v/%v", user, password, host, port, dbname),
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	var err error
	s.DB, err = sql.Open("pgx", s.dsn)
	if err != nil {
		return fmt.Errorf("cannot open pgx driver: %w", err)
	}

	return s.DB.PingContext(ctx)
}

func (s *Storage) Close(_ context.Context) error {
	return s.DB.Close()
}

func (s *Storage) Trucate(ctx context.Context) error {
	_, err := s.DB.ExecContext(ctx, "TRUNCATE slots CASCADE")
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) Migrate(_ context.Context, dirWithMigrations string) (err error) {
	//	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("cannot set dialect: %w", err)
	}

	if err := goose.Up(s.DB, dirWithMigrations); err != nil {
		return fmt.Errorf("cannot do up migration: %w", err)
	}

	return nil
}
