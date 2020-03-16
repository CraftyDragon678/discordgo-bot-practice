package main

import (
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

var messageLogs []string
var bots [2]*discordgo.Session

func runBot(token string, shardID int) (discord *discordgo.Session) {
	discord, err := discordgo.New("Bot " + token)
	checkErr("make bot", err)

	discord.ShardCount = len(bots)
	discord.ShardID = shardID

	discord.AddHandler(messageCreate)
	discord.AddHandlerOnce(func(s *discordgo.Session, r *discordgo.Ready) {
		s.UpdateStatus(0, "hello!")
		discord.
		// s.
	})

	err = discord.Open()
	checkErr("opening", err)

	return
}

func main() {
	bytes, err := ioutil.ReadFile("token")
	checkErr("load token file", err)

	for i := 0; i < len(bots); i++ {
		bots[i] = runBot(string(bytes), i)
		defer bots[i].Close()
	}

	e := echo.New()
	e.GET("/send", func(c echo.Context) error {
		return c.File("./send.html")
	})
	e.POST("/send", func(c echo.Context) error {
		sendAllChannel(c.FormValue("text"))
		return c.HTML(http.StatusOK, "<script>window.history.back()</script>")
	})
	e.GET("/log", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "<textarea>"+strings.Join(messageLogs, "\n")+"</textarea>")
	})

	go func() {
		e.Logger.Fatal(e.Start(":1232"))
	}()

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
					for i := 0; i < 6; i++ {
						go bot.ChannelMessageSend(channel.ID, content)
					}
					// break
				}
			}
		}
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
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
