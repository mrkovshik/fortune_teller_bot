package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	"github.com/mrkovshik/fortune_teller_bot/api/rest"
	"github.com/mrkovshik/fortune_teller_bot/internal/config"
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/bookstorage/local"
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/state_storage/inmemory"
	"github.com/mrkovshik/fortune_teller_bot/internal/update_processor/basic"
	"go.uber.org/zap"
)

func main() {
	_ = godotenv.Load()
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			logger.Fatal("can't sync logger")
		}
	}(logger)
	sugaredLogger := logger.Sugar()
	bookStorage := local.NewStorage(sugaredLogger)
	stateStorage := inmemory.NewStateStorage()
	commandProcessor := basic.NewUpdateProcessor(bookStorage, stateStorage, sugaredLogger)
	server := rest.NewRestAPIServer(commandProcessor, cfg, sugaredLogger)
	sugaredLogger.Fatal(server.RunServer(context.TODO()))
}
