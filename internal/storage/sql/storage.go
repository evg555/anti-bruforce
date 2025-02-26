package sqlstorage

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db *sqlx.DB
}

func New(host, port, user, password, dbname string) *Storage {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		user,
		password,
		host,
		port,
		dbname,
	)

	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		panic(fmt.Sprintf("database init error: %v", err))
	}

	return &Storage{
		db: db,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	err := s.db.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to db: %w", err)
	}

	return nil
}

func (s *Storage) Close(_ context.Context) error {
	return s.db.Close()
}
