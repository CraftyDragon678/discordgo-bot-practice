package cmd

import "bot-practice/framework"

// HelpCommand give help
func HelpCommand(ctx *framework.Context) {
	ctx.Reply("help!")
}
