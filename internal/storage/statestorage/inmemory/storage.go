package inmemory

import (
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/statestorage"
)

type StateStorage map[int64]*statestorage.ChatState // TODO: add mutex

func NewStateStorage() StateStorage { // TODO: добавить очистку по таймауту, чтобы не перегружать сторадж
	return make(StateStorage)
}

func (s StateStorage) Update(chatID int64, state *statestorage.ChatState) {
	s[chatID] = state
}

func (s StateStorage) Get(chatID int64) (*statestorage.ChatState, error) {
	state := s[chatID]
	return state, nil
}

func (s StateStorage) Add(chatID int64, state *statestorage.ChatState) {
	s[chatID] = state
}

func (s StateStorage) Clear(chatID int64) {
	delete(s, chatID)
}
