package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	"github.com/mrkovshik/fortune_teller_bot/api/rest"
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/state_storage/in_memory"
	"github.com/mrkovshik/fortune_teller_bot/internal/update_processor/basic"

	"github.com/mrkovshik/fortune_teller_bot/internal/config"
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/book_storage/local"

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
	defer logger.Sync()
	sugaredLogger := logger.Sugar()
	bookStorage := local.NewStorage(sugaredLogger)
	stateStorage := in_memory.NewStateStorage()
	commandProcessor := basic.NewUpdateProcessor(bookStorage, stateStorage, sugaredLogger)
	server := rest.NewRestAPIServer(commandProcessor, cfg, sugaredLogger)
	sugaredLogger.Fatal(server.RunServer(context.TODO()))
}
