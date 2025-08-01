package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

const botToken = "8310935771:AAFWTHsC4C-Yi1UKN22NQwVIvkosjrDdAao"

type Update struct {
	Message struct {
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
		Text string `json:"text"`
	} `json:"message"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	var update Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Println("Failed to decode update:", err)
		return
	}

	log.Printf("Got message: %s", update.Message.Text)

	sendMessage(update.Message.Chat.ID, "Вы написали: "+update.Message.Text)
}

func sendMessage(chatID int64, text string) {
	url := "https://api.telegram.org/bot" + botToken + "/sendMessage"
	payload := map[string]interface{}{
		"chat_id": chatID,
		"text":    text,
	}
	body, _ := json.Marshal(payload)
	http.Post(url, "application/json", bytes.NewBuffer(body))
}

func main() {
	http.HandleFunc("/", handler)
	port := os.Getenv("PORT") // Render требует переменную PORT
	if port == "" {
		port = "8080"
	}
	log.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
