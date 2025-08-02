package model

type Update struct {
	Message Message `json:"message"`
}
type Message struct {
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
}
type Chat struct {
	ID int64 `json:"id"`
}
