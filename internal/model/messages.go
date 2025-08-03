package model

type CallbackQuery struct {
	ID   string `json:"id"`
	From struct {
		ID int64 `json:"id"`
	} `json:"from"`
	Data string `json:"data"`
}

type Update struct {
	Message       *Message       `json:"message"`
	CallbackQuery *CallbackQuery `json:"callback_query" optional:"true"`
}
type Message struct {
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
}
type Chat struct {
	ID int64 `json:"id"`
}
