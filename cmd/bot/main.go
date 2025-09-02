package main

import (
	"context"
	"log"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mrkovshik/fortune_teller_bot/api/rest"
	"github.com/mrkovshik/fortune_teller_bot/internal/books-repository/embedded"
	"github.com/mrkovshik/fortune_teller_bot/internal/config"
	"github.com/mrkovshik/fortune_teller_bot/internal/poker"
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/books-meta/inmemory"
	postgresStorage "github.com/mrkovshik/fortune_teller_bot/internal/storage/books-meta/postgres"
	inmemorystep "github.com/mrkovshik/fortune_teller_bot/internal/storage/steps/inmemory"
	inmemoryuserdata "github.com/mrkovshik/fortune_teller_bot/internal/storage/userdata/inmemory"
	_ "github.com/mrkovshik/fortune_teller_bot/internal/templates"
	"github.com/mrkovshik/fortune_teller_bot/internal/updateprocessor"
	"github.com/mrkovshik/fortune_teller_bot/internal/updateprocessor/basic"
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
	var bookStorage updateprocessor.BookStorage
	bookStorage = inmemory.Storage
	if cfg.DatabaseURI != "" {
		bookStorage, err = postgresStorage.NewStorage(cfg)
		if err != nil {
			sugaredLogger.Fatal("can't init postgres storage", zap.Error(err))
		}
	}
	stateStorage := inmemorystep.NewStepStorage()
	userDataStorage := inmemoryuserdata.NewUserDataStorage()
	booksRepository := embedded.NewRepository()
	commandProcessor := basic.NewUpdateProcessor(bookStorage, booksRepository, stateStorage, userDataStorage, sugaredLogger, cfg)
	server := rest.NewRestAPIServer(commandProcessor, cfg, sugaredLogger)
	pokeTicker := time.NewTicker(time.Duration(cfg.PokingInterval) * time.Second)
	defer pokeTicker.Stop()
	urlPokerStopChanel := make(chan struct{})
	urlPoker := poker.NewPoker(sugaredLogger, cfg.PokingURL)
	go urlPoker.Poke(pokeTicker.C, urlPokerStopChanel)
	sugaredLogger.Fatal(server.RunServer(context.TODO()))
	<-urlPokerStopChanel
}
