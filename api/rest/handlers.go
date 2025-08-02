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
	telegramApiUrl = "https://api.telegram.org/bot"
	sendMessageUrl = "sendMessage"
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

		url := fmt.Sprintf("%s%s/%s", telegramApiUrl, s.cfg.Token, sendMessageUrl)

		reply, err := s.commandProcessor.ProcessCommand(update.Message.Text)
		if err != nil {
			s.logger.Error("ProcessCommand", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		s.logger.Infof("Sending reply: %s", reply)
		if err := sendMessage(update.Message.Chat.ID, reply, url); err != nil {
			s.logger.Error("sendMessage", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		c.AbortWithStatus(http.StatusOK)
	}
}

func sendMessage(chatID int64, text string, url string) error {
	payload := map[string]interface{}{
		"chat_id": chatID,
		"text":    text,
	}
	body, _ := json.Marshal(payload)
	_, err := http.Post(url, "application/json", bytes.NewBuffer(body)) // TODO: use lib
	return err
}
