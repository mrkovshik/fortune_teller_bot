package rest

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mrkovshik/fortune_teller_bot/internal/update_processor"

	"github.com/mrkovshik/fortune_teller_bot/internal/config"
	"github.com/mrkovshik/yandex_diploma/api"
	"go.uber.org/zap"
)

type restAPIServer struct {
	updateProcessor update_processor.UpdateProcessor
	logger          *zap.SugaredLogger
	cfg             *config.Config
}

func NewRestAPIServer(commandProcessor update_processor.UpdateProcessor, cfg *config.Config, logger *zap.SugaredLogger) api.Server {
	return &restAPIServer{
		updateProcessor: commandProcessor,
		logger:          logger,
		cfg:             cfg,
	}
}
func (s *restAPIServer) RunServer(ctx context.Context) error {
	router := gin.Default()
	router.POST("/telegram", s.MessageReplyHandler(ctx))
	router.POST("/", s.MessageReplyHandler(ctx))
	return router.Run(fmt.Sprintf("%s:%s", s.cfg.Host, s.cfg.Port))
}
