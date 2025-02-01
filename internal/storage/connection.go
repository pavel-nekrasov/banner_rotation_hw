package storage

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/stdlib" // need import pgx
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pavel-nekrasov/banner_rotation_hw/internal/config"
	goose "github.com/pressly/goose/v3"
)

type Connection struct {
	dsn string
	DB  *pgxpool.Pool
}

func NewConnection(config config.StorageConf) *Connection {
	return &Connection{
		dsn: fmt.Sprintf(
			"postgres://%v:%v@%v:%v/%v",
			config.User,
			config.Password,
			config.Host,
			config.Port,
			config.DBName,
		),
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
