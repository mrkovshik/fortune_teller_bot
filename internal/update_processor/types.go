package update_processor

import "github.com/mrkovshik/fortune_teller_bot/internal/model"

type UpdateProcessor interface {
	ProcessUpdate(update *model.Update) (map[string]interface{}, error)
}
