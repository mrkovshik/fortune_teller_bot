package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	"github.com/mrkovshik/fortune_teller_bot/api/rest"
	"github.com/mrkovshik/fortune_teller_bot/internal/command_processor/basic"
	"github.com/mrkovshik/fortune_teller_bot/internal/config"
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/local"
	"go.uber.org/zap"
)

var (
	sugaredLogger *zap.SugaredLogger
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
	sugaredLogger = logger.Sugar()

	storage := local.NewStorage(sugaredLogger)
	commandProcessor := basic.NewCommandProcessor(sugaredLogger, storage)
	server := rest.NewRestAPIServer(commandProcessor, cfg, sugaredLogger)
	sugaredLogger.Fatal(server.RunServer(context.TODO()))
}
