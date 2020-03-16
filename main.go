package main

import (
	"bot-practice/cmd"
	"bot-practice/framework"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/labstack/echo"
)

var (
	messageLogs []string
	bots        [2]*discordgo.Session

	cmdHandler *framework.CommandHandler
	prefix     string
	botID      string
)

func runBot(token string, shardID int) (discord *discordgo.Session) {
	discord, err := discordgo.New("Bot " + token)
	checkErr("make bot", err)

	discord.ShardCount = len(bots)
	discord.ShardID = shardID

	discord.AddHandler(commandHandler)
	discord.AddHandlerOnce(func(s *discordgo.Session, r *discordgo.Ready) {
		s.UpdateStatus(0, "hello!")
	})

	err = discord.Open()
	checkErr("opening", err)

	return
}

func runEcho() {
	e := echo.New()
	e.GET("/send", func(c echo.Context) error {
		return c.File("./send.html")
	})
	e.POST("/send", func(c echo.Context) error {
		sendAllChannel(c.FormValue("text"))
		return c.HTML(http.StatusOK, "<script>window.history.back()</script>")
	})
	e.POST("/embed", func(c echo.Context) error {
		sendAllChannel(c.FormValue("text"))
		return c.HTML(http.StatusOK, "<script>window.history.back()</script>")
	})
	e.GET("/log", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "<textarea>"+strings.Join(messageLogs, "\n")+"</textarea>")
	})

	go func() {
		e.Logger.Fatal(e.Start(":1232"))
	}()
}

func main() {
	bytes, err := ioutil.ReadFile("token")
	checkErr("load token file", err)

	cmdHandler = framework.NewCommandHandler()
	registerCommands()

	for i := 0; i < len(bots); i++ {
		bots[i] = runBot(string(bytes), i)
		defer bots[i].Close()
	}

	usr, err := bots[0].User("@me")
	checkErr("obtaining account details", err)

	botID = usr.ID

	runEcho()

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func sendAllChannel(content string) {
	for _, bot := range bots {
		guilds := bot.State.Guilds
		for _, guild := range guilds {
			for _, channel := range guild.Channels {
				perm, _ := bot.UserChannelPermissions(bot.State.User.ID, channel.ID)
				if channel.Type == discordgo.ChannelTypeGuildText && perm&discordgo.PermissionSendMessages != 0 {
					go bot.ChannelMessageSend(channel.ID, content)
					break
				}
			}
		}
	}
}

func commandHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	user := m.Author
	if user.ID == botID || user.Bot {
		return
	}
	content := m.Content
	args := strings.Fields(content)
	name := strings.ToLower(args[0])
	command, found := cmdHandler.Get(name)
	if found {
		channel, err := s.State.Channel(m.ChannelID)
		checkErr("getting channel", err)

		guild, err := s.State.Guild(channel.GuildID)
		checkErr("getting guild", err)

		ctx := framework.NewContext(s, guild, channel, user, message)
		ctx.Args = args[1:]
		c := *command
		c(ctx)
	}

	if m.Content == "ping" {
		for i := 1; i <= 9999999999; i++ {
		}
		s.ChannelMessageSend(m.ChannelID, "ping! "+s.HeartbeatLatency().String())
	}

	if m.Content == "pong" {
		for i := 0; i <= 5; i++ {
			s.ChannelMessageSend(m.ChannelID, "pong!")
		}
	}

	messageLogs = append(messageLogs, m.Content)
}

func checkErr(text string, err error) {
	if err != nil {
		for i := 0; i < len(bots); i++ {
			bots[i].Close()
		}
		log.Panicln("error occured in", text+"\n", err)
	}
}

func registerCommands() {
	cmdHandler.Register("help", cmd.HelpCommand, "help message!")
}
