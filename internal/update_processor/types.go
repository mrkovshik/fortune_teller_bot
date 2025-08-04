package update_processor

import "github.com/mrkovshik/fortune_teller_bot/internal/model"

type UpdateProcessor interface {
	ProcessMessage(message *model.Message) (map[string]interface{}, error)
	ProcessCallback(callback *model.CallbackQuery) (map[string]interface{}, error)
}
