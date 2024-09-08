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
	// load config
	ctx := context.Background()
	cfg, err := config.LoadConfig(".")
	if err != nil {
		logger.Error(ctx, "load config failed", err)
		return
	}

	// connect to postgres
	db, err := storage.InitPostgres(ctx, cfg.DatabaseDSN)
	if err != nil {
		panic(err)
	}

	// do migrations
	err = db.AutoMigrate(&users.User{}, &images.Image{}, &bookmarks.Bookmark{})
	if err != nil {
		logger.Error(ctx, "migration failed", err)
	}

	// set up app for listening
	a := app.NewApp(db)

	logger.Info(ctx, "starting server on port 8080")

	err = a.Listen(":8080")
	if err != nil {
		logger.Error(ctx, "start server failed", err)
		return
	}
}
