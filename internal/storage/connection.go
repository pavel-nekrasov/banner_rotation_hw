package storage

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/stdlib" // need import pgx
	"github.com/jackc/pgx/v4/pgxpool"
	goose "github.com/pressly/goose/v3"
)

type Connection struct {
	dsn string
	DB  *pgxpool.Pool
}

func NewConnection(host string, port int, dbname, user, password string) *Connection {
	return &Connection{
		dsn: fmt.Sprintf("postgres://%v:%v@%v:%v/%v", user, password, host, port, dbname),
	}
}

func (s *Connection) Connect(ctx context.Context) error {
	var err error

	s.DB, err = pgxpool.Connect(context.Background(), s.dsn)
	if err != nil {
		return fmt.Errorf("cannot open pgx driver: %w", err)
	}
	return s.DB.Ping(ctx)
}

func (s *Connection) Close() error {
	s.DB.Close()
	return nil
}

func (s *Connection) Migrate(_ context.Context, dirWithMigrations string) (err error) {
	//	goose.SetBaseFS(embedMigrations)
	db, err := sql.Open("pgx", s.dsn)
	if err != nil {
		return fmt.Errorf("cannot open pgx driver: %w", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("cannot set dialect: %w", err)
	}

	if err := goose.Up(db, dirWithMigrations); err != nil {
		return fmt.Errorf("cannot do up migration: %w", err)
	}

	return nil
}
