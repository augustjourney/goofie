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
		logger.Error(logger.Record{
			Message: "load config failed",
			Error:   err,
			Context: ctx,
		})
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
		logger.Error(logger.Record{
			Message: "migration failed",
			Error:   err,
			Context: ctx,
		})
	}

	// set up app for listening
	a := app.NewApp(db)

	logger.Info(logger.Record{
		Message: "starting server at port 8080",
		Context: ctx,
	})

	err = a.Listen(":8080")
	if err != nil {
		logger.Error(logger.Record{
			Message: "start server failed",
			Error:   err,
			Context: ctx,
		})
		return
	}
}
