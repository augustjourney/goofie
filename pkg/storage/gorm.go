package storage

import (
	"api/pkg/logger"
	"context"
	"github.com/jmoiron/sqlx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPostgres(ctx context.Context, dsn string) (*gorm.DB, error) {
	pg, err := sqlx.Open("pgx", dsn)
	if err != nil {
		logger.Error(logger.Record{
			Message: "[main InitPostgres] init postgres failed",
			Error:   err,
			Context: ctx,
		})
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
		logger.Error(logger.Record{
			Message: "[main postgres.Open] init postgres in postgres failed",
			Error:   err,
			Context: ctx,
		})
		return nil, err
	}

	logger.Info(logger.Record{
		Message: "[main postgres.Open] connected to postgres",
		Context: ctx,
	})

	return db, nil
}
