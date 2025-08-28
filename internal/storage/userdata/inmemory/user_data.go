package inmemory

import (
	"fmt"

	"github.com/mrkovshik/fortune_teller_bot/internal/storage/userdata"
)

type (
	UserDataStorage map[int64]userdata.UserData
)

func NewUserDataStorage() UserDataStorage {
	return make(UserDataStorage)
}

func (s UserDataStorage) GetUserData(chatID int64) (userdata.UserData, error) {
	userDataMap, exist := s[chatID]
	if !exist {
		return nil, fmt.Errorf("data for chatID %d not found", chatID)
	}

	return userDataMap, nil
}

func (s UserDataStorage) SaveUserData(chatID int64, key string, value string) error {
	_, exist := s[chatID]
	if !exist {
		s[chatID] = make(userdata.UserData)
	}
	s[chatID][key] = value
	return nil
}

func (s UserDataStorage) ClearUserData(chatID int64) error {
	delete(s, chatID)
	return nil
}
