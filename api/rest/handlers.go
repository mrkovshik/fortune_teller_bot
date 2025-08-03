package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrkovshik/fortune_teller_bot/internal/model"
)

const (
	telegramApiUrl    = "https://api.telegram.org/bot"
	sendMessageUrl    = "sendMessage"
	answerCallbackUrl = "answerCallbackQuery"
)

func (s *restAPIServer) MessageReplyHandler(ctx context.Context) func(c *gin.Context) {
	return func(c *gin.Context) {
		var update model.Update
		s.logger.Infof("Got request %s", c.Request.RequestURI)
		if c.Request.Body == nil {
			s.logger.Info("Empty body (maybe Telegram ping)")
			c.AbortWithStatus(http.StatusOK)
			return
		}
		if err := c.BindJSON(&update); err != nil {
			s.logger.Error("BindJSON", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		s.logger.Infof("Got message from chatID: %d : %s", update.Message.Chat.ID, update.Message.Text)

		reply, err := s.updateProcessor.ProcessUpdate(&update)
		if err != nil {
			s.logger.Error("ProcessUpdate", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		s.logger.Infof("Sending reply: %s", reply["text"])
		if err := s.sendMessage(reply); err != nil {
			s.logger.Error("sendMessage", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if update.CallbackQuery != nil {
			if err := s.answerCallbackQuery(update.CallbackQuery.ID); err != nil {
				s.logger.Error("sendMessage", err)
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
		}
		c.AbortWithStatus(http.StatusOK)
	}
}

func (s *restAPIServer) sendMessage(payload map[string]interface{}) error {
	url := fmt.Sprintf("%s%s/%s", telegramApiUrl, s.cfg.Token, sendMessageUrl)
	body, _ := json.Marshal(payload)
	_, err := http.Post(url, "application/json", bytes.NewBuffer(body)) // TODO: use lib
	return err
}

func (s *restAPIServer) answerCallbackQuery(callbackID string) error {
	payload := map[string]interface{}{
		"callback_query_id": callbackID,
	}
	url := fmt.Sprintf("%s%s/%s", telegramApiUrl, s.cfg.Token, answerCallbackUrl)
	body, _ := json.Marshal(payload)
	_, err := http.Post(url, "application/json", bytes.NewBuffer(body)) // TODO: use lib
	return err
}
