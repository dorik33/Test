package store

import (
	"context"
	"fmt"

	"github.com/dorik33/Test/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type Store struct {
	pool           *pgxpool.Pool
	logger         *logrus.Logger
	SongRepository *SongRepository
}

func NewConnection(cfg *config.Config, logger *logrus.Logger) (*Store, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	)

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	store := &Store{
		pool:   pool,
		logger: logger,
	}

	store.SongRepository = &SongRepository{store: store}

	return store, nil
}

func (s *Store) Close() {
	s.pool.Close()
}
