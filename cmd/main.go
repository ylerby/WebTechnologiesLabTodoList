package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	application "backend/internal/app"
	"backend/internal/cache"
	"backend/internal/database"
	"backend/internal/handlers"
	zapLogger "backend/internal/logger"
	"go.uber.org/zap"
)

const (
	appLoggerKey   = "component"
	appLoggerValue = "todo_list"
	MongoUriKey    = "MONGO_URI"
)

func main() {
	logger, err := zapLogger.New()
	if err != nil {
		log.Fatalf("ошибка инициализации логгера: %s", err)
	}

	defer func() {
		_ = logger.Sync()
	}()

	appLogger := logger.With(zap.String(appLoggerKey, appLoggerValue))
	appLogger.Info("инициализация логгера")

	MongoURI := os.Getenv(MongoUriKey)
	appCache, err := cache.New(MongoURI)
	if err != nil {
		appLogger.Errorf("ошибка при подключении к mongoDB - %s", err)
		return
	}

	appLogger.Infof("инициализация кеша - %v", appCache)

	appDatabase := database.New()
	err = appDatabase.Connect()
	if err != nil {
		appLogger.Errorf("ошибка при подключении к postgreSQL - %s", err)
		return
	}

	appLogger.Infof("инициализация базы данных - %v", appDatabase)

	appHandlers := handlers.New(appDatabase, appCache, appLogger)
	app := application.New(appHandlers.InitRoutes())

	go func() {
		_ = app.Run()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	err = app.Stop()
	if err != nil {
		appLogger.Fatalf("ошибка при завершении работы сервера - %s", err)
	}
}
