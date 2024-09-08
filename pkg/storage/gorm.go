package storage

import (
	"api/pkg/logger"
	"context"
	"github.com/jmoiron/sqlx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitPostgres creates a new connection to postgres using gorm
func InitPostgres(ctx context.Context, dsn string) (*gorm.DB, error) {
	pg, err := sqlx.Open("pgx", dsn)
	if err != nil {
		logger.Error(ctx, "init postgres failed", err)
		return nil, err
	}

	err = pg.Ping()
	if err != nil {
		logger.Error(ctx, "init postgres failed", err)
		return nil, err
	}

	gormPostgresConn := postgres.New(postgres.Config{
		Conn: pg,
	})
	gormConfig := &gorm.Config{
		Logger: nil,
	}
	db, err := gorm.Open(gormPostgresConn, gormConfig)
	if err != nil {
		logger.Error(ctx, "open gorm postgres failed", err)
		return nil, err
	}

	logger.Info(ctx, "connected to postgres")

	return db, nil
}
