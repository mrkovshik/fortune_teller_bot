package rdb

import (
	"context"

	"github.com/mrkovshik/fortune_teller_bot/internal/config"
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/userdata"
	"github.com/redis/go-redis/v9"
)

type (
	UserDataStorage struct {
		rdb *redis.Client
	}
)

func NewUserDataStorage(ctx context.Context, cfg *config.RDB) (*UserDataStorage, error) {
	opt := &redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	client := redis.NewClient(opt)
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &UserDataStorage{
		rdb: client,
	}, nil
}

func (s UserDataStorage) GetBookID(ctx context.Context, chatID int64) (int64, error) {
	return s.rdb.Get(ctx, userdata.GenerateKey(chatID, userdata.BookIDKey)).Int64()
}

func (s UserDataStorage) GetLanguage(ctx context.Context, chatID int64) (config.Language, error) {
	return config.Language(s.rdb.Get(ctx, userdata.GenerateKey(chatID, userdata.LanguageKey)).String()), nil
}

func (s UserDataStorage) SaveBookID(ctx context.Context, chatID, bookID int64) error {
	return s.rdb.Set(ctx, userdata.GenerateKey(chatID, userdata.BookIDKey), bookID, 0).Err()
}

func (s UserDataStorage) SaveLanguage(ctx context.Context, chatID int64, language config.Language) error {
	return s.rdb.Set(ctx, userdata.GenerateKey(chatID, userdata.LanguageKey), language, 0).Err()
}
