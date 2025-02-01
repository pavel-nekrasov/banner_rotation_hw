package storage

import (
	"context"

	_ "github.com/jackc/pgx/stdlib" // need import pgx
)

type Storage struct {
	conn *Connection
}

func NewStorage(conn *Connection) *Storage {
	return &Storage{
		conn: conn,
	}
}

func (s *Storage) Trucate(ctx context.Context) error {
	_, err := s.conn.DB.Exec(ctx, "TRUNCATE slots CASCADE")
	if err != nil {
		return err
	}

	return nil
}
