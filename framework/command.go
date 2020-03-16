package framework

type (
	// Command executes
	Command func(*Context)

	commandStruct struct {
		command Command
		help    string
	}

	// CmdMap list of command
	CmdMap map[string]commandStruct

	// CommandHandler handler of command
	CommandHandler struct {
		cmds CmdMap
	}
)

// NewCommandHandler makes new command handler
func NewCommandHandler() *CommandHandler {
	return &CommandHandler{make(CmdMap)}
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
	cmdstruct := commandStruct{command: command, help: helpmsg}
	handler.cmds[name] = cmdstruct

	// if len(name) > 1 {
	// 	handler.cmds[name[:1]] = cmdstruct
	// }
}

// GetHelp returns help of command
func (command commandStruct) GetHelp() string {
	return command.help
}
