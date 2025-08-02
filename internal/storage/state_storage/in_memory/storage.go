package in_memory

import (
	"errors"

	"github.com/mrkovshik/fortune_teller_bot/internal/model"
)

type StateStorage map[int64]model.ChatState // TODO: add mutex

func NewStateStorage() StateStorage { // TODO: добавить очистку по таймауту, чтобы не перегружать сторадж
	return make(StateStorage)
}

func (s StateStorage) Update(chatID int64, state model.ChatState) {
	s[chatID] = state
}

func (s StateStorage) Get(chatID int64) (model.ChatState, error) {
	state, ok := s[chatID]
	if !ok {
		return model.ChatState{}, errors.New("chat does not exist")
	}
	return state, nil
}

func (s StateStorage) Add(chatID int64, state model.ChatState) {
	s[chatID] = state
}

func (s StateStorage) Clear(chatID int64) {
	delete(s, chatID)
}
