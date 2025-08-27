package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrkovshik/fortune_teller_bot/internal/model"
	"go.uber.org/zap"
)

const (
	telegramAPIURL    = "https://api.telegram.org/bot"
	sendMessageURL    = "sendMessage"
	answerCallbackURL = "answerCallbackQuery"
)

func (s *restAPIServer) MessageReplyHandler(_ context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if !c.Writer.Written() {
				c.String(http.StatusOK, "OK")
			}
		}()

		var update model.Update

		s.logger.Infof("Got request %s", c.Request.RequestURI)
		if c.Request.Body == nil {
			s.logger.Info("Empty body (maybe Telegram ping)")
			return
		}

		if err := c.ShouldBindJSON(&update); err != nil {
			s.logger.Warn("bad JSON", zap.Error(err))
			return
		}

		switch {
		case update.Message != nil:
			s.logger.Infof("Got message from chatID: %d : %s", update.Message.Chat.ID, update.Message.Text)

			reply, err := s.updateProcessor.ProcessMessage(update.Message)
			if err != nil {
				s.logger.Warn("ProcessMessage", zap.Error(err))
				if err := s.sendMessage(map[string]interface{}{
					"chat_id": update.Message.Chat.ID,
					"text":    "⚠️ Что-то пошло не так. Попробуйте ещё раз позже.",
				}); err != nil {
					s.logger.Warn("sendMessage", zap.Error(err))
				}
				return
			}
			if err := s.sendMessage(reply); err != nil {
				s.logger.Warn("sendMessage", zap.Error(err))
				return
			}

		case update.CallbackQuery != nil:
			s.logger.Infof("Got callback from chatID: %d", update.CallbackQuery.From.ID)

			reply, err := s.updateProcessor.ProcessCallback(update.CallbackQuery)
			if err != nil {
				s.logger.Warn("ProcessCallback", zap.Error(err))
				_ = s.answerCallbackQuery(update.CallbackQuery.ID)
				return
			}
			if err := s.sendMessage(reply); err != nil {
				s.logger.Warn("sendMessage", zap.Error(err))
				_ = s.answerCallbackQuery(update.CallbackQuery.ID)
				return
			}
			_ = s.answerCallbackQuery(update.CallbackQuery.ID)

		default:
			s.logger.Info("Unsupported update type")
			return
		}
	}
}

func (s *restAPIServer) sendMessage(payload map[string]interface{}) error {
	url := fmt.Sprintf("%s%s/%s", telegramAPIURL, s.cfg.Token, sendMessageURL)
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body)) // TODO: use lib
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	s.logger.Infof("Telegram response: %s", string(respBody))
	return nil
}

func (s *restAPIServer) answerCallbackQuery(callbackID string) error {
	payload := map[string]interface{}{
		"callback_query_id": callbackID,
	}
	url := fmt.Sprintf("%s%s/%s", telegramAPIURL, s.cfg.Token, answerCallbackURL)
	body, _ := json.Marshal(payload)
	_, err := http.Post(url, "application/json", bytes.NewBuffer(body)) // TODO: use lib
	return err
}
