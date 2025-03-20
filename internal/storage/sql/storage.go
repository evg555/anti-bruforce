package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/evg555/antibrutforce/internal/storage"
	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver is required for database/sql package
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

func (s *Storage) Save(ctx context.Context, subnet storage.Subnet, listType string) error {
	query := fmt.Sprintf("INSERT INTO %s (subnet) VALUES ($1)", listType)

	_, err := s.db.ExecContext(ctx, query, subnet.Address)
	if err != nil {
		return fmt.Errorf("failed to save subnet %s to %s: %w", subnet.Address, listType, err)
	}

	return nil
}

func (s *Storage) Find(ctx context.Context, address, listType string) (*storage.Subnet, error) {
	var rows []*storage.Subnet

	query := fmt.Sprintf("SELECT * FROM %s WHERE subnet=$1 LIMIT 1", listType)

	err := s.db.SelectContext(ctx, &rows, query, address)
	if err != nil {
		return nil, fmt.Errorf("failed to find subnet %s in %s: %w", address, listType, err)
	}

	return rows[0], nil
}

func (s *Storage) Delete(ctx context.Context, address, listType string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE subnet=$1", listType)

	_, err := s.db.ExecContext(ctx, query, address)
	if err != nil {
		return fmt.Errorf("failed to delete subnet %s in %s: %w", address, listType, err)
	}

	return nil
}

func (s *Storage) IsInList(ctx context.Context, address, listType string) (bool, error) {
	var res int

	query := fmt.Sprintf("SELECT 1 FROM %s WHERE $1::inet << subnet", listType)

	row := s.db.QueryRowxContext(ctx, query, address)

	err := row.Scan(&res)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check ip %s in %s: %w", address, listType, err)
	}

	return res == 1, nil
}
