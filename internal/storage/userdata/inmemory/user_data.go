package inmemory

import "fmt"

type (
	UserData        map[string]string
	UserDataStorage map[int64]UserData
)

func NewUserDataStorage() UserDataStorage {
	return make(UserDataStorage)
}

func (s UserDataStorage) GetUserData(chatID int64, key string) (string, error) {
	userDataMap, exist := s[chatID]
	if !exist {
		return "", fmt.Errorf("data for chatID %d not found", chatID)
	}
	value, exist := userDataMap[key]
	if !exist {
		return "", fmt.Errorf("data for chatID %d for key %s not found", chatID, key)
	}
	return value, nil
}

func (s UserDataStorage) ClearUserData(chatID int64) error {
	delete(s, chatID)
	return nil
}
