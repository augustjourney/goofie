package main

import (
	app "api/internal"
	"api/internal/bookmarks"
	"api/internal/images"
	"api/internal/users"
	"api/pkg/config"
	"api/pkg/logger"
	"api/pkg/storage"
	"context"
)

func main() {
	ctx := context.Background()
	cfg, err := config.LoadConfig(".")
	if err != nil {
		logger.Error(logger.Record{
			Message: "[main LoadConfig] load config failed",
			Error:   err,
			Context: ctx,
		})
		return
	}
	db, err := storage.InitPostgres(ctx, cfg.DatabaseDSN)
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&users.User{}, &images.Image{}, &bookmarks.Bookmark{})
	if err != nil {
		logger.Error(logger.Record{
			Message: "[main db.AutoMigrate] migration failed",
			Error:   err,
			Context: ctx,
		})
	}

	a := app.NewApp(db)

	logger.Info(logger.Record{
		Message: "[main app.Listen] starting server at port 8080",
		Context: ctx,
	})

	err = a.Listen(":8080")
	if err != nil {
		logger.Error(logger.Record{
			Message: "[main app.Listen] start server failed",
			Error:   err,
			Context: ctx,
		})
		return
	}
}
