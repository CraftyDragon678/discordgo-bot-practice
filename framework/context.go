package framework

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Context of data
type Context struct {
	Discord      *discordgo.Session
	Guild        *discordgo.Guild
	VoiceChannel *discordgo.Channel
	TextChannel  *discordgo.Channel
	User         *discordgo.User
	Message      *discordgo.MessageCreate
	Args         []string
}

// NewContext returns new context
func NewContext(discord *discordgo.Session, guild *discordgo.Guild,
	textChannel *discordgo.Channel, user *discordgo.User,
	message *discordgo.MessageCreate) *Context {
	ctx := new(Context)
	ctx.Discord = discord
	ctx.Guild = guild
	ctx.TextChannel = textChannel
	ctx.User = user
	ctx.Message = message

	return ctx
}

// Reply to channel
func (ctx Context) Reply(context string) *discordgo.Message {
	msg, err := ctx.Discord.ChannelMessageSend(ctx.TextChannel.ID, context)
	if err != nil {
		fmt.Println("Error ", err)
		return nil
	}
	return msg
}

// GetVoiceChannel returns current voice channel if bot has joined voice channel, otherwise returns nil
func (ctx *Context) GetVoiceChannel() *discordgo.Channel {
	if ctx.VoiceChannel != nil {
		return ctx.VoiceChannel
	}
	for _, state := range ctx.Guild.VoiceStates {
		if state.UserID == ctx.User.ID {
			channel, _ := ctx.Discord.Channel(state.ChannelID)
			ctx.VoiceChannel = channel
			return channel
		}
	}
	return nil
}
