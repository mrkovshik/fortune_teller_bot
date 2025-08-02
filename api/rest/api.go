package rest

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mrkovshik/fortune_teller_bot/internal/command_processor"
	"github.com/mrkovshik/fortune_teller_bot/internal/config"
	"github.com/mrkovshik/yandex_diploma/api"
	"go.uber.org/zap"
)

type restAPIServer struct {
	commandProcessor command_processor.CommandProcessor
	logger           *zap.SugaredLogger
	cfg              *config.Config
}

func NewRestAPIServer(commandProcessor command_processor.CommandProcessor, cfg *config.Config, logger *zap.SugaredLogger) api.Server {
	return &restAPIServer{
		commandProcessor: commandProcessor,
		logger:           logger,
		cfg:              cfg,
	}
}
func (s *restAPIServer) RunServer(ctx context.Context) error {
	router := gin.Default()
	router.POST("/telegram", s.MessageReplyHandler(ctx))
	router.POST("/", s.MessageReplyHandler(ctx))
	return router.Run(fmt.Sprintf("%s:%s", s.cfg.Host, s.cfg.Port))
}
