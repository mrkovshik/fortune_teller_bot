package command_processor

type CommandProcessor interface {
	ProcessCommand(command string) (string, error)
}
