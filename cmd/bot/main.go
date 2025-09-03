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
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/userdata/rdb"
	_ "github.com/mrkovshik/fortune_teller_bot/internal/templates"
	"github.com/mrkovshik/fortune_teller_bot/internal/updateprocessor"
	"github.com/mrkovshik/fortune_teller_bot/internal/updateprocessor/basic"
	"go.uber.org/zap"
)

func main() {
	ctx := context.TODO()
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

	if cfg.DatabaseURI != "" {
		bookStorage, err = postgresStorage.NewStorage(cfg)
		if err != nil {
			sugaredLogger.Fatal("can't init postgres books meta storage", zap.Error(err))
		}
		sugaredLogger.Debug("init postgres books meta storage")
	} else {
		bookStorage = inmemory.Storage
		sugaredLogger.Debug("init inmemory books meta storage")
	}
	stateStorage := inmemorystep.NewStepStorage()
	var userDataStorage updateprocessor.UserDataStorage

	if cfg.RDB != nil {
		userDataStorage, err = rdb.NewUserDataStorage(ctx, cfg.RDB)
		if err != nil {
			sugaredLogger.Fatal("can't init rdb user data storage", zap.Error(err))
		}
		sugaredLogger.Debug("init rdb user data storage")
	} else {
		userDataStorage = inmemoryuserdata.NewUserDataStorage()
		sugaredLogger.Debug("init inmemory user data storage")
	}
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
