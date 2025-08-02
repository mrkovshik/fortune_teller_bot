package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/mrkovshik/fortune_teller_bot/internal/command_processor/basic"
	"go.uber.org/zap"
)

type Update struct {
	Message struct {
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
		Text string `json:"text"`
	} `json:"message"`
}

var (
	url           string
	sugaredLogger *zap.SugaredLogger
)

const (
	telegramApiUrl = "https://api.telegram.org/bot"
	sendMessageUrl = "sendMessage"
)

func handler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Failed to read body:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(body) == 0 {
		log.Println("Empty body (maybe Telegram ping)")
		w.WriteHeader(http.StatusOK)
		return
	}

	var update Update
	if err := json.Unmarshal(body, &update); err != nil {
		log.Println("Failed to decode JSON:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sugaredLogger.Infof("Got message: %s", update.Message.Text)
	commandProcessor := basic.NewCommandProcessor(sugaredLogger)
	message, err := commandProcessor.ProcessCommand(update.Message.Text)
	if err != nil {
		sugaredLogger.Infof("Failed to process message: %s", err)
	}
	sendMessage(update.Message.Chat.ID, message)
	w.WriteHeader(http.StatusOK)
}

func sendMessage(chatID int64, text string) {
	payload := map[string]interface{}{
		"chat_id": chatID,
		"text":    text,
	}
	body, _ := json.Marshal(payload)
	http.Post(url, "application/json", bytes.NewBuffer(body))
}

func main() {
	_ = godotenv.Load()
	http.HandleFunc("/", handler)
	token := os.Getenv("TELEGRAM_TOKEN")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()
	sugaredLogger = logger.Sugar()
	url = fmt.Sprintf("%s/%s/%s", telegramApiUrl, token, sendMessageUrl)
	log.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
