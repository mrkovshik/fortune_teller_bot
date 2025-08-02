package basic

import (
	"fmt"
	"strings"

	"github.com/mrkovshik/fortune_teller_bot/internal/command_processor"
	store "github.com/mrkovshik/fortune_teller_bot/internal/storage"
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/local"
	"go.uber.org/zap"
)

type CommandProcessor struct {
	logger  *zap.SugaredLogger
	storage store.Storage
}

func NewCommandProcessor(logger *zap.SugaredLogger, storage store.Storage) *CommandProcessor {
	return &CommandProcessor{
		logger:  logger,
		storage: storage,
	}
}

func (cp *CommandProcessor) ProcessCommand(command string) (string, error) {
	switch command {
	case command_processor.ListBooksCommandName:
		list, err := cp.storage.ListBooks()
		if err != nil {
			return "", fmt.Errorf(`failed to list books: %w`, err)
		}
		return strings.Join(list, "/n"), nil // TODO: use template
	case command_processor.GetMagicCommandName:
		return cp.storage.GetRandomSentenceFromBook(local.DorianGreyTitle)
	default:
		return "Неизвестная команда", nil
	}
}
