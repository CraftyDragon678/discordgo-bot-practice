package framework

type (
	Command func(Context)

	CommandStruct struct {
		command Command
		help    string
	}

	cmdMap map[string]CommandStruct

	// CommandHandler handler of command
	CommandHandler struct {
		cmds cmdMap
	}
)

// NewCommandHandler makes new command handler
func NewCommandHandler() *CommandHandler {
	return &CommandHandler(make(CmdMap))
}

// GetCmds returns cmds of handler
func (handler CommandHandler) GetCmds() CmdMap {
	return handler.cmds
}

// Get returns cmd by name
func (handler CommandHandler) Get(name string) (*Command, bool) {
	cmd, found := handler.cmds[name]
	return &cmd.command, found
}

// Register command
func (handler CommandHandler) Register(name string, command Command, helpmsg string) {
	cmdstruct := CommandStruct{command: command, help: helpmsg}
	handler.cmds[name] = cmdstruct

	// if len(name) > 1 {
	// 	handler.cmds[name[:1]] = cmdstruct
	// }
}

// GetHelp returns help of command
func (command CommandStruct) GetHelp() string {
	return command.help
}